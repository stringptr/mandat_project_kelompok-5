# ETL Logging & Audit Trail - User Guide

## 📋 Overview

Sistem ETL logging menyediakan audit trail lengkap untuk setiap eksekusi ETL, termasuk:
- ✅ Tracking waktu mulai dan selesai
- ✅ Jumlah rows yang diproses (insert, update, delete)
- ✅ Status eksekusi (SUCCESS/FAILED)
- ✅ Error message detail
- ✅ Step-by-step execution log
- ✅ Performance metrics
- ✅ Incremental loading metadata

---

## 🗄️ Database Schema

### **Tabel Utama:**

#### 1. **etl_log** - Main execution log
```sql
- id_etl_log (PK)
- etl_name (e.g., 'MasterToWarehouse')
- pipeline_name (e.g., 'dim_pasien')
- execution_id (unique identifier)
- start_time, end_time, duration_seconds
- status ('RUNNING', 'SUCCESS', 'FAILED')
- rows_read, rows_inserted, rows_updated, rows_deleted
- error_message, error_code
- parameters (JSONB - untuk LAST_ETL_RUN, dll)
```

#### 2. **etl_log_detail** - Step-by-step log
```sql
- id_etl_log_detail (PK)
- id_etl_log (FK)
- step_name (e.g., 'Read Pasien Data', 'Execute UPSERT')
- start_time, end_time, duration_ms
- status, rows_input, rows_output, rows_error
- error_message
```

#### 3. **etl_metadata** - Incremental loading metadata
```sql
- etl_name (unique)
- last_successful_run (untuk incremental filter)
- high_watermark
- total_runs, successful_runs, failed_runs
- total_rows_processed
- avg_duration_seconds
```

---

## 🚀 Cara Penggunaan

### **1. Di Awal Pipeline - Start Execution**

**Tambahkan step di awal setiap pipeline:**

**Step Name:** `00_Start_ETL_Log`
**Transform Type:** Execute SQL
**Connection:** Warehouse

**SQL:**
```sql
SELECT public.start_etl_execution(
    'MasterToWarehouse',              -- ETL name
    'dim_pasien',                      -- Pipeline name
    'RunAll',                          -- Workflow name (optional)
    'pasien',                          -- Source table (optional)
    'DIM_PASIEN',                      -- Target table (optional)
    '{"LAST_ETL_RUN": "${LAST_ETL_RUN}"}'::JSONB  -- Parameters (optional)
) as execution_id;
```

**Output:** Simpan `execution_id` ke variable untuk digunakan di step selanjutnya.

---

### **2. Di Akhir Pipeline - End Execution (Success)**

**Step Name:** `99_End_ETL_Log_Success`
**Transform Type:** Execute SQL
**Connection:** Warehouse

**SQL:**
```sql
SELECT public.end_etl_execution_success(
    '${execution_id}',    -- Execution ID dari step 1
    100,                   -- Rows inserted
    50,                    -- Rows updated
    0                      -- Rows deleted
);
```

**Note:** Rows count bisa di-hardcode atau ambil dari counter transform.

---

### **3. Di Error Handler - End Execution (Failed)**

**Step Name:** `99_End_ETL_Log_Failed`
**Transform Type:** Execute SQL (in Error Handling)
**Connection:** Warehouse

**SQL:**
```sql
SELECT public.end_etl_execution_failed(
    '${execution_id}',
    'Error occurred during ETL execution',  -- Error message
    'ETL_001'                                -- Error code (optional)
);
```

---

### **4. Optional: Log Individual Steps**

**Untuk logging step-by-step execution:**

**Step Name:** `Log_Transform_Step`
**Transform Type:** Execute SQL
**Connection:** Warehouse

**SQL:**
```sql
SELECT public.log_etl_step(
    '${execution_id}',
    'Read Pasien Data',     -- Step name
    'TableInput',            -- Step type
    1000,                    -- Rows input
    950,                     -- Rows output
    'SUCCESS',               -- Status
    NULL                     -- Error message (if any)
);
```

---

## 📊 Query untuk Monitoring & Audit

### **1. Lihat Eksekusi Terakhir**

```sql
SELECT 
    etl_name,
    pipeline_name,
    execution_id,
    start_time,
    end_time,
    duration_seconds,
    status,
    rows_total,
    error_message
FROM etl_log
ORDER BY start_time DESC
LIMIT 10;
```

### **2. Lihat Summary per ETL**

```sql
SELECT * FROM v_etl_execution_summary
WHERE etl_name = 'MasterToWarehouse'
ORDER BY start_time DESC;
```

### **3. Performance Metrics**

```sql
SELECT * FROM v_etl_performance
ORDER BY avg_duration_seconds DESC;
```

### **4. Recent Failures**

```sql
SELECT * FROM v_etl_recent_failures
ORDER BY start_time DESC;
```

### **5. Success Rate by Pipeline**

```sql
SELECT 
    pipeline_name,
    COUNT(*) as total_runs,
    SUM(CASE WHEN status = 'SUCCESS' THEN 1 ELSE 0 END) as successful,
    SUM(CASE WHEN status = 'FAILED' THEN 1 ELSE 0 END) as failed,
    ROUND(
        SUM(CASE WHEN status = 'SUCCESS' THEN 1 ELSE 0 END)::NUMERIC / 
        COUNT(*) * 100, 2
    ) as success_rate_percent
FROM etl_log
WHERE start_time >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY pipeline_name
ORDER BY success_rate_percent ASC;
```

### **6. ETL Duration Trend**

```sql
SELECT 
    DATE_TRUNC('day', start_time) as execution_date,
    pipeline_name,
    AVG(duration_seconds) as avg_duration,
    MIN(duration_seconds) as min_duration,
    MAX(duration_seconds) as max_duration,
    COUNT(*) as execution_count
FROM etl_log
WHERE status = 'SUCCESS'
  AND start_time >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY DATE_TRUNC('day', start_time), pipeline_name
ORDER BY execution_date DESC, pipeline_name;
```

### **7. Rows Processed Over Time**

```sql
SELECT 
    DATE_TRUNC('day', start_time) as execution_date,
    SUM(rows_inserted) as total_inserted,
    SUM(rows_updated) as total_updated,
    SUM(rows_deleted) as total_deleted,
    SUM(rows_total) as total_processed
FROM etl_log
WHERE status = 'SUCCESS'
  AND start_time >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY DATE_TRUNC('day', start_time)
ORDER BY execution_date DESC;
```

### **8. Failed Steps Detail**

```sql
SELECT 
    el.pipeline_name,
    el.execution_id,
    el.start_time,
    eld.step_name,
    eld.error_message,
    eld.rows_input,
    eld.rows_output,
    eld.rows_error
FROM etl_log el
JOIN etl_log_detail eld ON el.id_etl_log = eld.id_etl_log
WHERE eld.status = 'FAILED'
  AND el.start_time >= CURRENT_DATE - INTERVAL '7 days'
ORDER BY el.start_time DESC;
```

---

## 🔍 Incremental Loading dengan ETL Metadata

### **Cara Kerja:**

1. **Sebelum ETL:**
```sql
SELECT public.get_last_etl_run('dim_pasien');
-- Returns: '2026-06-15 10:30:00'
```

2. **Gunakan di WHERE clause:**
```sql
WHERE updated_at >= public.get_last_etl_run('dim_pasien')
  AND is_deleted = false
```

3. **Setelah ETL Success:**
```sql
-- Otomatis di-update oleh end_etl_execution_success()
-- last_successful_run = CURRENT_TIMESTAMP
```

### **Manual Update (Jika Perlu):**

```sql
UPDATE etl_metadata 
SET last_successful_run = '2026-06-15 12:00:00'
WHERE etl_name = 'dim_pasien';
```

---

## 🛠️ Maintenance

### **Cleanup Old Logs (> 90 days)**

```sql
DELETE FROM etl_log
WHERE start_time < CURRENT_DATE - INTERVAL '90 days'
  AND status IN ('SUCCESS', 'FAILED');
```

### **Archive Logs to History Table**

```sql
-- Create history table first
CREATE TABLE etl_log_history (LIKE etl_log INCLUDING ALL);

-- Move old logs
INSERT INTO etl_log_history
SELECT * FROM etl_log
WHERE start_time < CURRENT_DATE - INTERVAL '90 days';

-- Delete from main table
DELETE FROM etl_log
WHERE start_time < CURRENT_DATE - INTERVAL '90 days';
```

### **Reset ETL Metadata (Full Reload)**

```sql
UPDATE etl_metadata
SET last_successful_run = '1900-01-01'::TIMESTAMPTZ,
    last_watermark = '1900-01-01'::TIMESTAMPTZ
WHERE etl_name = 'dim_pasien';
```

---

## 📈 Dashboard Metrics (untuk BI Tools)

### **KPI Queries untuk Dashboard:**

**1. Today's ETL Status**
```sql
SELECT 
    COUNT(*) FILTER (WHERE status = 'SUCCESS') as successful,
    COUNT(*) FILTER (WHERE status = 'FAILED') as failed,
    COUNT(*) FILTER (WHERE status = 'RUNNING') as running,
    SUM(rows_total) as total_rows_processed
FROM etl_log
WHERE start_time >= CURRENT_DATE;
```

**2. Average ETL Duration (Last 7 Days)**
```sql
SELECT 
    pipeline_name,
    ROUND(AVG(duration_seconds), 2) as avg_seconds,
    COUNT(*) as runs
FROM etl_log
WHERE start_time >= CURRENT_DATE - INTERVAL '7 days'
  AND status = 'SUCCESS'
GROUP BY pipeline_name;
```

**3. Error Rate by Pipeline**
```sql
SELECT 
    pipeline_name,
    COUNT(*) FILTER (WHERE status = 'FAILED') as errors,
    COUNT(*) as total,
    ROUND(
        COUNT(*) FILTER (WHERE status = 'FAILED')::NUMERIC / 
        COUNT(*) * 100, 2
    ) as error_rate
FROM etl_log
WHERE start_time >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY pipeline_name
HAVING COUNT(*) FILTER (WHERE status = 'FAILED') > 0
ORDER BY error_rate DESC;
```

---

## ✅ Best Practices

1. **Always log execution start and end**
   - Gunakan `start_etl_execution()` di awal
   - Gunakan `end_etl_execution_success/failed()` di akhir

2. **Include meaningful error messages**
   - Capture full error details
   - Include step name where error occurred

3. **Monitor regularly**
   - Setup alerts untuk failed ETL
   - Review performance metrics weekly
   - Cleanup old logs periodically

4. **Use parameters JSON**
   - Store LAST_ETL_RUN dan parameter lain
   - Berguna untuk troubleshooting

5. **Track row counts accurately**
   - Gunakan counter transform di Hop
   - Validasi dengan target table counts

---

## 🚨 Alerting & Notifications

### **Query for Alert System:**

```sql
-- ETL yang fail dalam 1 jam terakhir
SELECT 
    etl_name,
    pipeline_name,
    execution_id,
    error_message
FROM etl_log
WHERE status = 'FAILED'
  AND start_time >= CURRENT_TIMESTAMP - INTERVAL '1 hour';
```

### **Email Alert Template:**

```
Subject: ⚠️ ETL Failure Alert - {pipeline_name}

Pipeline: {pipeline_name}
Execution ID: {execution_id}
Start Time: {start_time}
Error: {error_message}

Check details at: http://hop-gui/logs/{execution_id}
```

---

## 📚 References

- **ETL Logging Best Practices:** https://docs.apache.org/hop/latest/
- **PostgreSQL JSONB:** https://www.postgresql.org/docs/current/datatype-json.html
- **Audit Trail Standards:** ISO 27001, GDPR compliance

---

**Created:** 2026-06-16  
**Version:** 1.0  
**Maintainer:** Data Engineering Team
