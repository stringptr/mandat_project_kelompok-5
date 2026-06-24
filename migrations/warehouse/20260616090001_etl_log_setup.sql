-- Create helper functions
CREATE OR REPLACE FUNCTION public.start_etl_execution(
    p_etl_name VARCHAR, p_pipeline_name VARCHAR,
    p_workflow_name VARCHAR DEFAULT 'N/A', p_source_table VARCHAR DEFAULT 'N/A',
    p_target_table VARCHAR DEFAULT 'N/A', p_parameters JSONB DEFAULT '{}'::jsonb
) RETURNS VARCHAR AS $$
DECLARE
    v_execution_id VARCHAR(50);
    v_id BIGINT;
BEGIN
    v_execution_id := 'ETL_' || TO_CHAR(CURRENT_TIMESTAMP, 'YYYYMMDDHH24MISS') || '_' || LPAD(FLOOR(RANDOM() * 10000)::TEXT, 4, '0');
    INSERT INTO public.etl_log (
        etl_name, pipeline_name, workflow_name, execution_id, start_time, status, source_table, target_table, parameters,
        source_system, target_system, hop_server, hop_version, user_name, hostname
    )
    VALUES (
        p_etl_name, p_pipeline_name, COALESCE(p_workflow_name, 'N/A'), v_execution_id, CURRENT_TIMESTAMP, 'RUNNING', COALESCE(p_source_table, 'N/A'), COALESCE(p_target_table, 'N/A'), COALESCE(p_parameters, '{}'::jsonb),
        COALESCE(p_parameters->>'SOURCE_SYSTEM', 'OLTP_SiGizi'),
        COALESCE(p_parameters->>'TARGET_SYSTEM', 'Data_Warehouse'),
        COALESCE(p_parameters->>'HOP_SERVER', 'Local_Hop'),
        COALESCE(p_parameters->>'HOP_VERSION', 'Unknown'),
        CURRENT_USER,
        COALESCE(inet_client_addr()::TEXT, 'localhost')
    )
    RETURNING id_etl_log INTO v_id;
    INSERT INTO public.etl_metadata (etl_name, status, last_run_start, table_name)
    VALUES (p_etl_name, 'RUNNING', CURRENT_TIMESTAMP, COALESCE(p_target_table, 'N/A'))
    ON CONFLICT (etl_name) DO UPDATE SET status='RUNNING', last_run_start=CURRENT_TIMESTAMP, total_runs=etl_metadata.total_runs+1, updated_at=CURRENT_TIMESTAMP;
    RETURN v_execution_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION public.end_etl_execution_success(
    p_execution_id VARCHAR, p_rows_inserted BIGINT DEFAULT 0,
    p_rows_updated BIGINT DEFAULT 0, p_rows_deleted BIGINT DEFAULT 0
) RETURNS VOID AS $$
DECLARE
    v_start_time TIMESTAMPTZ;
    v_duration INTEGER;
    v_etl_name VARCHAR(100);
BEGIN
    SELECT start_time, etl_name INTO v_start_time, v_etl_name FROM public.etl_log WHERE execution_id = p_execution_id;
    v_duration := EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - v_start_time))::INTEGER;
    UPDATE public.etl_log SET end_time=CURRENT_TIMESTAMP, duration_seconds=v_duration, status='SUCCESS',
        rows_inserted=p_rows_inserted, rows_updated=p_rows_updated, rows_deleted=p_rows_deleted, updated_at=CURRENT_TIMESTAMP
    WHERE execution_id = p_execution_id;
    UPDATE public.etl_metadata SET last_successful_run=CURRENT_TIMESTAMP, status='IDLE',
        successful_runs=successful_runs+1, total_rows_processed=total_rows_processed+(p_rows_inserted+p_rows_updated+p_rows_deleted),
        last_run_end=CURRENT_TIMESTAMP, last_run_status='SUCCESS', last_run_rows=p_rows_inserted+p_rows_updated+p_rows_deleted,
        updated_at=CURRENT_TIMESTAMP
    WHERE etl_name = v_etl_name;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION public.end_etl_execution_failed(
    p_execution_id VARCHAR, p_error_message TEXT, p_error_code VARCHAR DEFAULT '-'
) RETURNS VOID AS $$
DECLARE
    v_start_time TIMESTAMPTZ;
    v_duration INTEGER;
    v_etl_name VARCHAR(100);
BEGIN
    SELECT start_time, etl_name INTO v_start_time, v_etl_name FROM public.etl_log WHERE execution_id = p_execution_id;
    v_duration := EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - v_start_time))::INTEGER;
    UPDATE public.etl_log SET end_time=CURRENT_TIMESTAMP, duration_seconds=v_duration, status='FAILED',
        error_message=COALESCE(p_error_message, '-'), error_code=COALESCE(p_error_code, '-'), updated_at=CURRENT_TIMESTAMP
    WHERE execution_id = p_execution_id;
    UPDATE public.etl_metadata SET status='FAILED', failed_runs=failed_runs+1,
        last_run_end=CURRENT_TIMESTAMP, last_run_status='FAILED', last_error_message=COALESCE(p_error_message, '-'), updated_at=CURRENT_TIMESTAMP
    WHERE etl_name = v_etl_name;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION public.get_last_etl_run(p_etl_name VARCHAR) RETURNS TIMESTAMPTZ AS $$
DECLARE v_last_run TIMESTAMPTZ;
BEGIN
    SELECT COALESCE(last_successful_run, '1900-01-01'::TIMESTAMPTZ) INTO v_last_run FROM public.etl_metadata WHERE etl_name = p_etl_name;
    IF v_last_run IS NULL THEN RETURN '1900-01-01'::TIMESTAMPTZ; END IF;
    RETURN v_last_run;
END;
$$ LANGUAGE plpgsql;

-- ============================================================
-- INSERT SAMPLE LOG DATA - Simulasi 1 kali RunAll berhasil
-- ============================================================

DO $$
DECLARE
    t0  TIMESTAMPTZ := NOW() - INTERVAL '5 minutes';
    ex0 VARCHAR := 'ETL_20260623184500_1000';
    ex1 VARCHAR := 'ETL_20260623184501_1001';
    ex2 VARCHAR := 'ETL_20260623184501_1002';
    ex3 VARCHAR := 'ETL_20260623184501_1003';
    ex4 VARCHAR := 'ETL_20260623184504_1004';
    ex5 VARCHAR := 'ETL_20260623184507_1005';
    ex6 VARCHAR := 'ETL_20260623184513_1006';
    ex7 VARCHAR := 'ETL_20260623184513_1007';
    ex8 VARCHAR := 'ETL_20260623184518_1008';
    ex9 VARCHAR := 'ETL_20260623184518_1009';
    id0 BIGINT; id1 BIGINT; id2 BIGINT; id3 BIGINT; id4 BIGINT;
    id5 BIGINT; id6 BIGINT; id7 BIGINT; id8 BIGINT; id9 BIGINT;
BEGIN

-- RunAll Workflow (parent)
INSERT INTO etl_log(etl_name,pipeline_name,workflow_name,execution_id,start_time,end_time,duration_seconds,status,rows_inserted,rows_updated,rows_deleted,source_table,target_table,parameters)
VALUES('MasterToWarehouse','RunAll_Workflow','RunAll',ex0,t0+INTERVAL'0s',t0+INTERVAL'42s',42,'SUCCESS',0,0,0,'Multiple_Master_Tables','All_DIM_and_FACT_Tables',
  '{"LAST_ETL_RUN":"2026-06-22 10:00:00","pipelines":["dim_waktu","dim_lokasi","dim_petugas","dim_posyandu","dim_pasien","dim_anak","dim_ibu_hamil","fact_pemeriksaan","fact_imunisasi"]}'::JSONB)
RETURNING id_etl_log INTO id0;

-- dim_waktu
INSERT INTO etl_log(etl_name,pipeline_name,workflow_name,execution_id,start_time,end_time,duration_seconds,status,rows_inserted,rows_updated,rows_deleted,source_table,target_table,parameters)
VALUES('MasterToWarehouse','dim_waktu','RunAll',ex1,t0+INTERVAL'1s',t0+INTERVAL'4s',3,'SUCCESS',365,0,0,'N/A','DIM_WAKTU','{"LAST_ETL_RUN":"2026-06-22 10:00:00"}'::JSONB)
RETURNING id_etl_log INTO id1;

-- dim_lokasi
INSERT INTO etl_log(etl_name,pipeline_name,workflow_name,execution_id,start_time,end_time,duration_seconds,status,rows_inserted,rows_updated,rows_deleted,source_table,target_table,parameters)
VALUES('MasterToWarehouse','dim_lokasi','RunAll',ex2,t0+INTERVAL'1s',t0+INTERVAL'3s',2,'SUCCESS',5,2,0,'lokasi','DIM_LOKASI','{"LAST_ETL_RUN":"2026-06-22 10:00:00"}'::JSONB)
RETURNING id_etl_log INTO id2;

-- dim_petugas
INSERT INTO etl_log(etl_name,pipeline_name,workflow_name,execution_id,start_time,end_time,duration_seconds,status,rows_inserted,rows_updated,rows_deleted,source_table,target_table,parameters)
VALUES('MasterToWarehouse','dim_petugas','RunAll',ex3,t0+INTERVAL'1s',t0+INTERVAL'3s',2,'SUCCESS',3,1,0,'user_account','DIM_PETUGAS','{"LAST_ETL_RUN":"2026-06-22 10:00:00"}'::JSONB)
RETURNING id_etl_log INTO id3;

-- dim_posyandu
INSERT INTO etl_log(etl_name,pipeline_name,workflow_name,execution_id,start_time,end_time,duration_seconds,status,rows_inserted,rows_updated,rows_deleted,source_table,target_table,parameters)
VALUES('MasterToWarehouse','dim_posyandu','RunAll',ex4,t0+INTERVAL'4s',t0+INTERVAL'6s',2,'SUCCESS',2,0,0,'posyandu','DIM_POSYANDU','{"LAST_ETL_RUN":"2026-06-22 10:00:00"}'::JSONB)
RETURNING id_etl_log INTO id4;

-- dim_pasien
INSERT INTO etl_log(etl_name,pipeline_name,workflow_name,execution_id,start_time,end_time,duration_seconds,status,rows_inserted,rows_updated,rows_deleted,source_table,target_table,parameters)
VALUES('MasterToWarehouse','dim_pasien','RunAll',ex5,t0+INTERVAL'7s',t0+INTERVAL'12s',5,'SUCCESS',12,3,0,'pasien','DIM_PASIEN','{"LAST_ETL_RUN":"2026-06-22 10:00:00"}'::JSONB)
RETURNING id_etl_log INTO id5;

-- dim_anak
INSERT INTO etl_log(etl_name,pipeline_name,workflow_name,execution_id,start_time,end_time,duration_seconds,status,rows_inserted,rows_updated,rows_deleted,source_table,target_table,parameters)
VALUES('MasterToWarehouse','dim_anak','RunAll',ex6,t0+INTERVAL'13s',t0+INTERVAL'16s',3,'SUCCESS',8,2,0,'anak','DIM_ANAK','{"LAST_ETL_RUN":"2026-06-22 10:00:00"}'::JSONB)
RETURNING id_etl_log INTO id6;

-- dim_ibu_hamil
INSERT INTO etl_log(etl_name,pipeline_name,workflow_name,execution_id,start_time,end_time,duration_seconds,status,rows_inserted,rows_updated,rows_deleted,source_table,target_table,parameters)
VALUES('MasterToWarehouse','dim_ibu_hamil','RunAll',ex7,t0+INTERVAL'13s',t0+INTERVAL'17s',4,'SUCCESS',4,1,0,'ibu_hamil','DIM_IBU_HAMIL','{"LAST_ETL_RUN":"2026-06-22 10:00:00"}'::JSONB)
RETURNING id_etl_log INTO id7;

-- fact_pemeriksaan
INSERT INTO etl_log(etl_name,pipeline_name,workflow_name,execution_id,start_time,end_time,duration_seconds,status,rows_inserted,rows_updated,rows_deleted,source_table,target_table,parameters)
VALUES('MasterToWarehouse','fact_pemeriksaan','RunAll',ex8,t0+INTERVAL'18s',t0+INTERVAL'35s',17,'SUCCESS',87,0,0,'hasil_pemeriksaan','FACT_PEMERIKSAAN','{"LAST_ETL_RUN":"2026-06-22 10:00:00"}'::JSONB)
RETURNING id_etl_log INTO id8;

-- fact_imunisasi
INSERT INTO etl_log(etl_name,pipeline_name,workflow_name,execution_id,start_time,end_time,duration_seconds,status,rows_inserted,rows_updated,rows_deleted,source_table,target_table,parameters)
VALUES('MasterToWarehouse','fact_imunisasi','RunAll',ex9,t0+INTERVAL'18s',t0+INTERVAL'30s',12,'SUCCESS',45,0,0,'jadwal_imunisasi','FACT_IMUNISASI','{"LAST_ETL_RUN":"2026-06-22 10:00:00"}'::JSONB)
RETURNING id_etl_log INTO id9;

-- ============================================================
-- etl_log_detail untuk dim_pasien (id5)
-- ============================================================
INSERT INTO etl_log_detail(id_etl_log,step_name,step_type,step_order,start_time,end_time,duration_ms,status,rows_input,rows_output,rows_error)
VALUES
(id5,'00_Start_ETL_Log','TableInput',1,t0+INTERVAL'7.000s',t0+INTERVAL'7.120s',120,'SUCCESS',0,1,0),
(id5,'Read Pasien with Calculations','TableInput',2,t0+INTERVAL'7.120s',t0+INTERVAL'9.450s',2330,'SUCCESS',0,15,0),
(id5,'Build SQL','ScriptValueMod',3,t0+INTERVAL'9.450s',t0+INTERVAL'9.780s',330,'SUCCESS',15,15,0),
(id5,'Execute UPSERT','ExecSqlRow',4,t0+INTERVAL'9.780s',t0+INTERVAL'11.900s',2120,'SUCCESS',15,15,0),
(id5,'99_End_ETL_Log_Success','ExecSql',5,t0+INTERVAL'11.900s',t0+INTERVAL'12.010s',110,'SUCCESS',0,0,0);

-- ============================================================
-- etl_log_detail untuk fact_pemeriksaan (id8)
-- ============================================================
INSERT INTO etl_log_detail(id_etl_log,step_name,step_type,step_order,start_time,end_time,duration_ms,status,rows_input,rows_output,rows_error)
VALUES
(id8,'00_Start_ETL_Log','TableInput',1,t0+INTERVAL'18.000s',t0+INTERVAL'18.130s',130,'SUCCESS',0,1,0),
(id8,'Read Pemeriksaan Data','TableInput',2,t0+INTERVAL'18.130s',t0+INTERVAL'22.500s',4370,'SUCCESS',0,87,0),
(id8,'Build Flags & SQL','ScriptValueMod',3,t0+INTERVAL'22.500s',t0+INTERVAL'23.100s',600,'SUCCESS',87,87,0),
(id8,'Execute INSERT FACT','ExecSqlRow',4,t0+INTERVAL'23.100s',t0+INTERVAL'34.800s',11700,'SUCCESS',87,87,0),
(id8,'99_End_ETL_Log_Success','ExecSql',5,t0+INTERVAL'34.800s',t0+INTERVAL'35.000s',200,'SUCCESS',0,0,0);

-- ============================================================
-- etl_metadata
-- ============================================================
INSERT INTO etl_metadata(etl_name,table_name,last_successful_run,status,total_runs,successful_runs,failed_runs,total_rows_processed,avg_duration_seconds,last_run_start,last_run_end,last_run_status,last_run_rows)
VALUES
('dim_waktu','DIM_WAKTU',NOW(),'IDLE',3,3,0,365*3,3.2,t0+INTERVAL'1s',t0+INTERVAL'4s','SUCCESS',365),
('dim_lokasi','DIM_LOKASI',NOW(),'IDLE',3,3,0,21,2.1,t0+INTERVAL'1s',t0+INTERVAL'3s','SUCCESS',7),
('dim_petugas','DIM_PETUGAS',NOW(),'IDLE',3,3,0,12,2.0,t0+INTERVAL'1s',t0+INTERVAL'3s','SUCCESS',4),
('dim_posyandu','DIM_POSYANDU',NOW(),'IDLE',3,3,0,6,1.8,t0+INTERVAL'4s',t0+INTERVAL'6s','SUCCESS',2),
('dim_pasien','DIM_PASIEN',NOW(),'IDLE',3,2,1,45,5.0,t0+INTERVAL'7s',t0+INTERVAL'12s','SUCCESS',15),
('dim_anak','DIM_ANAK',NOW(),'IDLE',3,3,0,30,3.1,t0+INTERVAL'13s',t0+INTERVAL'16s','SUCCESS',10),
('dim_ibu_hamil','DIM_IBU_HAMIL',NOW(),'IDLE',3,3,0,15,4.0,t0+INTERVAL'13s',t0+INTERVAL'17s','SUCCESS',5),
('fact_pemeriksaan','FACT_PEMERIKSAAN',NOW(),'IDLE',3,3,0,261,16.8,t0+INTERVAL'18s',t0+INTERVAL'35s','SUCCESS',87),
('fact_imunisasi','FACT_IMUNISASI',NOW(),'IDLE',3,3,0,135,11.5,t0+INTERVAL'18s',t0+INTERVAL'30s','SUCCESS',45)
ON CONFLICT (etl_name) DO NOTHING;

END $$;
