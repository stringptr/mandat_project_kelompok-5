-- +goose Up
-- Migration: Create ETL logging system for audit trail
-- Created: 2026-06-16 09:00:00

-- =====================================================
-- ETL_LOG Table: Main logging table
-- =====================================================
CREATE TABLE IF NOT EXISTS public.etl_log (
    id_etl_log BIGSERIAL PRIMARY KEY,
    
    -- ETL Identification
    etl_name VARCHAR(100) NOT NULL,
    pipeline_name VARCHAR(200) NOT NULL,
    workflow_name VARCHAR(200),
    
    -- Execution Tracking
    execution_id VARCHAR(50) NOT NULL UNIQUE,
    parent_execution_id VARCHAR(50),
    
    -- Timestamps
    start_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMPTZ,
    duration_seconds INTEGER,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'RUNNING',
    -- Status values: 'RUNNING', 'SUCCESS', 'FAILED', 'CANCELLED'
    
    -- Data Metrics
    rows_read BIGINT DEFAULT 0,
    rows_inserted BIGINT DEFAULT 0,
    rows_updated BIGINT DEFAULT 0,
    rows_deleted BIGINT DEFAULT 0,
    rows_rejected BIGINT DEFAULT 0,
    rows_total BIGINT GENERATED ALWAYS AS (rows_inserted + rows_updated + rows_deleted) STORED,
    
    -- Source & Target
    source_system VARCHAR(100),
    source_table VARCHAR(200),
    target_system VARCHAR(100),
    target_table VARCHAR(200),
    
    -- Error Handling
    error_message TEXT,
    error_code VARCHAR(50),
    error_line INTEGER,
    
    -- Performance
    peak_memory_mb NUMERIC(10,2),
    cpu_time_seconds NUMERIC(10,2),
    
    -- Metadata
    hop_server VARCHAR(100),
    hop_version VARCHAR(50),
    user_name VARCHAR(100),
    hostname VARCHAR(200),
    
    -- Additional Info
    parameters JSONB,
    notes TEXT,
    
    -- Audit
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_etl_log_etl_name ON public.etl_log(etl_name);
CREATE INDEX idx_etl_log_pipeline_name ON public.etl_log(pipeline_name);
CREATE INDEX idx_etl_log_execution_id ON public.etl_log(execution_id);
CREATE INDEX idx_etl_log_start_time ON public.etl_log(start_time DESC);
CREATE INDEX idx_etl_log_status ON public.etl_log(status);
CREATE INDEX idx_etl_log_etl_status ON public.etl_log(etl_name, status);

-- =====================================================
-- ETL_LOG_DETAIL Table: Step-by-step logging
-- =====================================================
CREATE TABLE IF NOT EXISTS public.etl_log_detail (
    id_etl_log_detail BIGSERIAL PRIMARY KEY,
    id_etl_log BIGINT NOT NULL REFERENCES public.etl_log(id_etl_log) ON DELETE CASCADE,
    
    -- Step Information
    step_name VARCHAR(200) NOT NULL,
    step_type VARCHAR(50),
    step_order INTEGER,
    
    -- Timing
    start_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMPTZ,
    duration_ms INTEGER,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'RUNNING',
    
    -- Metrics
    rows_input BIGINT DEFAULT 0,
    rows_output BIGINT DEFAULT 0,
    rows_error BIGINT DEFAULT 0,
    
    -- Error Details
    error_message TEXT,
    error_description TEXT,
    
    -- Additional Info
    sql_executed TEXT,
    parameters JSONB,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_etl_log_detail_etl_log ON public.etl_log_detail(id_etl_log);
CREATE INDEX idx_etl_log_detail_step_name ON public.etl_log_detail(step_name);
CREATE INDEX idx_etl_log_detail_status ON public.etl_log_detail(status);

-- =====================================================
-- ETL_METADATA Table: For incremental loading
-- =====================================================
CREATE TABLE IF NOT EXISTS public.etl_metadata (
    id_etl_metadata SERIAL PRIMARY KEY,
    
    etl_name VARCHAR(100) NOT NULL UNIQUE,
    table_name VARCHAR(200),
    
    -- Incremental Tracking
    last_successful_run TIMESTAMPTZ,
    last_watermark TIMESTAMPTZ,
    high_watermark TIMESTAMPTZ,
    
    -- Status
    status VARCHAR(20) DEFAULT 'IDLE',
    last_execution_id VARCHAR(50),
    
    -- Statistics
    total_runs BIGINT DEFAULT 0,
    successful_runs BIGINT DEFAULT 0,
    failed_runs BIGINT DEFAULT 0,
    total_rows_processed BIGINT DEFAULT 0,
    avg_duration_seconds NUMERIC(10,2),
    
    -- Last Run Info
    last_run_start TIMESTAMPTZ,
    last_run_end TIMESTAMPTZ,
    last_run_status VARCHAR(20),
    last_run_rows BIGINT,
    last_error_message TEXT,
    
    -- Audit
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_etl_metadata_etl_name ON public.etl_metadata(etl_name);
CREATE INDEX idx_etl_metadata_status ON public.etl_metadata(status);

-- =====================================================
-- Helper Functions
-- =====================================================

-- Function: Start ETL Execution
CREATE OR REPLACE FUNCTION public.start_etl_execution(
    p_etl_name VARCHAR,
    p_pipeline_name VARCHAR,
    p_workflow_name VARCHAR DEFAULT NULL,
    p_source_table VARCHAR DEFAULT NULL,
    p_target_table VARCHAR DEFAULT NULL,
    p_parameters JSONB DEFAULT NULL
)
RETURNS VARCHAR AS $$
DECLARE
    v_execution_id VARCHAR(50);
    v_id_etl_log BIGINT;
BEGIN
    -- Generate unique execution ID
    v_execution_id := 'ETL_' || TO_CHAR(CURRENT_TIMESTAMP, 'YYYYMMDDHH24MISS') || '_' || 
                      LPAD(FLOOR(RANDOM() * 10000)::TEXT, 4, '0');
    
    -- Insert log record
    INSERT INTO public.etl_log (
        etl_name, pipeline_name, workflow_name, execution_id,
        start_time, status, source_table, target_table, parameters
    ) VALUES (
        p_etl_name, p_pipeline_name, p_workflow_name, v_execution_id,
        CURRENT_TIMESTAMP, 'RUNNING', p_source_table, p_target_table, p_parameters
    ) RETURNING id_etl_log INTO v_id_etl_log;
    
    -- Update metadata status
    INSERT INTO public.etl_metadata (etl_name, status, last_execution_id, last_run_start, table_name)
    VALUES (p_etl_name, 'RUNNING', v_execution_id, CURRENT_TIMESTAMP, p_target_table)
    ON CONFLICT (etl_name) DO UPDATE SET
        status = 'RUNNING',
        last_execution_id = v_execution_id,
        last_run_start = CURRENT_TIMESTAMP,
        total_runs = etl_metadata.total_runs + 1,
        updated_at = CURRENT_TIMESTAMP;
    
    RETURN v_execution_id;
END;
$$ LANGUAGE plpgsql;

-- Function: End ETL Execution (Success)
CREATE OR REPLACE FUNCTION public.end_etl_execution_success(
    p_execution_id VARCHAR,
    p_rows_inserted BIGINT DEFAULT 0,
    p_rows_updated BIGINT DEFAULT 0,
    p_rows_deleted BIGINT DEFAULT 0
)
RETURNS VOID AS $$
DECLARE
    v_start_time TIMESTAMPTZ;
    v_duration INTEGER;
    v_etl_name VARCHAR(100);
BEGIN
    -- Get start time and calculate duration
    SELECT start_time, etl_name INTO v_start_time, v_etl_name
    FROM public.etl_log
    WHERE execution_id = p_execution_id;
    
    v_duration := EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - v_start_time))::INTEGER;
    
    -- Update log record
    UPDATE public.etl_log SET
        end_time = CURRENT_TIMESTAMP,
        duration_seconds = v_duration,
        status = 'SUCCESS',
        rows_inserted = p_rows_inserted,
        rows_updated = p_rows_updated,
        rows_deleted = p_rows_deleted,
        updated_at = CURRENT_TIMESTAMP
    WHERE execution_id = p_execution_id;
    
    -- Update metadata
    UPDATE public.etl_metadata SET
        last_successful_run = CURRENT_TIMESTAMP,
        last_watermark = CURRENT_TIMESTAMP,
        status = 'IDLE',
        successful_runs = successful_runs + 1,
        total_rows_processed = total_rows_processed + (p_rows_inserted + p_rows_updated + p_rows_deleted),
        last_run_end = CURRENT_TIMESTAMP,
        last_run_status = 'SUCCESS',
        last_run_rows = p_rows_inserted + p_rows_updated + p_rows_deleted,
        avg_duration_seconds = (COALESCE(avg_duration_seconds, 0) * (successful_runs) + v_duration) / (successful_runs + 1),
        updated_at = CURRENT_TIMESTAMP
    WHERE etl_name = v_etl_name;
END;
$$ LANGUAGE plpgsql;

-- Function: End ETL Execution (Failed)
CREATE OR REPLACE FUNCTION public.end_etl_execution_failed(
    p_execution_id VARCHAR,
    p_error_message TEXT,
    p_error_code VARCHAR DEFAULT NULL
)
RETURNS VOID AS $$
DECLARE
    v_start_time TIMESTAMPTZ;
    v_duration INTEGER;
    v_etl_name VARCHAR(100);
BEGIN
    SELECT start_time, etl_name INTO v_start_time, v_etl_name
    FROM public.etl_log
    WHERE execution_id = p_execution_id;
    
    v_duration := EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - v_start_time))::INTEGER;
    
    UPDATE public.etl_log SET
        end_time = CURRENT_TIMESTAMP,
        duration_seconds = v_duration,
        status = 'FAILED',
        error_message = p_error_message,
        error_code = p_error_code,
        updated_at = CURRENT_TIMESTAMP
    WHERE execution_id = p_execution_id;
    
    UPDATE public.etl_metadata SET
        status = 'FAILED',
        failed_runs = failed_runs + 1,
        last_run_end = CURRENT_TIMESTAMP,
        last_run_status = 'FAILED',
        last_error_message = p_error_message,
        updated_at = CURRENT_TIMESTAMP
    WHERE etl_name = v_etl_name;
END;
$$ LANGUAGE plpgsql;

-- Function: Get Last Successful Run Timestamp
CREATE OR REPLACE FUNCTION public.get_last_etl_run(p_etl_name VARCHAR)
RETURNS TIMESTAMPTZ AS $$
DECLARE
    v_last_run TIMESTAMPTZ;
BEGIN
    SELECT COALESCE(last_successful_run, '1900-01-01'::TIMESTAMPTZ)
    INTO v_last_run
    FROM public.etl_metadata
    WHERE etl_name = p_etl_name;
    
    IF v_last_run IS NULL THEN
        RETURN '1900-01-01'::TIMESTAMPTZ;
    END IF;
    
    RETURN v_last_run;
END;
$$ LANGUAGE plpgsql;

-- Function: Log ETL Step
CREATE OR REPLACE FUNCTION public.log_etl_step(
    p_execution_id VARCHAR,
    p_step_name VARCHAR,
    p_step_type VARCHAR DEFAULT NULL,
    p_rows_input BIGINT DEFAULT 0,
    p_rows_output BIGINT DEFAULT 0,
    p_status VARCHAR DEFAULT 'SUCCESS',
    p_error_message TEXT DEFAULT NULL
)
RETURNS VOID AS $$
DECLARE
    v_id_etl_log BIGINT;
BEGIN
    SELECT id_etl_log INTO v_id_etl_log
    FROM public.etl_log
    WHERE execution_id = p_execution_id;
    
    INSERT INTO public.etl_log_detail (
        id_etl_log, step_name, step_type, start_time, end_time,
        status, rows_input, rows_output, error_message
    ) VALUES (
        v_id_etl_log, p_step_name, p_step_type, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP,
        p_status, p_rows_input, p_rows_output, p_error_message
    );
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- Views for Reporting
-- =====================================================

-- View: ETL Execution Summary
CREATE OR REPLACE VIEW public.v_etl_execution_summary AS
SELECT 
    el.etl_name,
    el.pipeline_name,
    el.execution_id,
    el.start_time,
    el.end_time,
    el.duration_seconds,
    el.status,
    el.rows_total,
    el.rows_inserted,
    el.rows_updated,
    el.error_message,
    COUNT(eld.id_etl_log_detail) as total_steps,
    SUM(CASE WHEN eld.status = 'FAILED' THEN 1 ELSE 0 END) as failed_steps
FROM public.etl_log el
LEFT JOIN public.etl_log_detail eld ON el.id_etl_log = eld.id_etl_log
GROUP BY el.id_etl_log, el.etl_name, el.pipeline_name, el.execution_id, 
         el.start_time, el.end_time, el.duration_seconds, el.status, 
         el.rows_total, el.rows_inserted, el.rows_updated, el.error_message;

-- View: ETL Performance Metrics
CREATE OR REPLACE VIEW public.v_etl_performance AS
SELECT 
    etl_name,
    table_name,
    total_runs,
    successful_runs,
    failed_runs,
    ROUND(successful_runs::NUMERIC / NULLIF(total_runs, 0) * 100, 2) as success_rate,
    total_rows_processed,
    avg_duration_seconds,
    last_successful_run,
    last_run_status,
    last_error_message
FROM public.etl_metadata
ORDER BY total_runs DESC;

-- View: Recent ETL Failures
CREATE OR REPLACE VIEW public.v_etl_recent_failures AS
SELECT 
    el.etl_name,
    el.pipeline_name,
    el.execution_id,
    el.start_time,
    el.end_time,
    el.error_message,
    el.error_code,
    STRING_AGG(eld.step_name || ': ' || eld.error_message, '; ') as step_errors
FROM public.etl_log el
LEFT JOIN public.etl_log_detail eld ON el.id_etl_log = eld.id_etl_log AND eld.status = 'FAILED'
WHERE el.status = 'FAILED'
  AND el.start_time >= CURRENT_TIMESTAMP - INTERVAL '7 days'
GROUP BY el.id_etl_log, el.etl_name, el.pipeline_name, el.execution_id,
         el.start_time, el.end_time, el.error_message, el.error_code
ORDER BY el.start_time DESC;

-- =====================================================
-- Comments for Documentation
-- =====================================================

COMMENT ON TABLE public.etl_log IS 'Main ETL execution log for audit trail';
COMMENT ON TABLE public.etl_log_detail IS 'Step-by-step ETL execution details';
COMMENT ON TABLE public.etl_metadata IS 'ETL metadata for incremental loading and statistics';

COMMENT ON COLUMN public.etl_log.execution_id IS 'Unique identifier for each ETL run';
COMMENT ON COLUMN public.etl_log.rows_total IS 'Computed column: sum of inserted + updated + deleted';
COMMENT ON COLUMN public.etl_log.parameters IS 'JSON object containing ETL parameters (e.g., LAST_ETL_RUN)';

COMMENT ON FUNCTION public.start_etl_execution IS 'Initialize ETL execution and return execution_id';
COMMENT ON FUNCTION public.end_etl_execution_success IS 'Mark ETL execution as successful with metrics';
COMMENT ON FUNCTION public.end_etl_execution_failed IS 'Mark ETL execution as failed with error details';
COMMENT ON FUNCTION public.get_last_etl_run IS 'Get last successful run timestamp for incremental loading';
COMMENT ON FUNCTION public.log_etl_step IS 'Log individual ETL step execution';

-- +goose Down
DROP VIEW IF EXISTS public.v_etl_recent_failures;
DROP VIEW IF EXISTS public.v_etl_performance;
DROP VIEW IF EXISTS public.v_etl_execution_summary;

DROP FUNCTION IF EXISTS public.log_etl_step;
DROP FUNCTION IF EXISTS public.get_last_etl_run;
DROP FUNCTION IF EXISTS public.end_etl_execution_failed;
DROP FUNCTION IF EXISTS public.end_etl_execution_success;
DROP FUNCTION IF EXISTS public.start_etl_execution;

DROP TABLE IF EXISTS public.etl_log_detail CASCADE;
DROP TABLE IF EXISTS public.etl_metadata CASCADE;
DROP TABLE IF EXISTS public.etl_log CASCADE;
