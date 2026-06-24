# Requirements Document

## Introduction

Sistem ETL SiGizi Data Warehouse saat ini memiliki 9 pipeline MasterToWarehouse (7 dimensi + 2 fakta) yang beroperasi dalam mode incremental loading, namun belum memiliki audit trail lengkap untuk compliance dan troubleshooting. Fitur ini akan mengintegrasikan ETL logging ke semua pipeline MasterToWarehouse dengan memanfaatkan fungsi database yang sudah tersedia (start_etl_execution, end_etl_execution_success, end_etl_execution_failed) untuk mencatat setiap eksekusi pipeline ke tabel etl_log di data warehouse.

## Glossary

- **ETL_Pipeline**: Apache Hop pipeline (.hpl file) yang melakukan Extract, Transform, Load dari Master ke Warehouse
- **Warehouse**: PostgreSQL database yang menyimpan data warehouse dan tabel etl_log
- **Master**: PostgreSQL database sumber yang menyimpan data OLTP
- **Start_ETL_Log_Step**: Transform step Execute SQL yang memanggil start_etl_execution() di awal pipeline
- **End_ETL_Log_Success_Step**: Transform step Execute SQL yang memanggil end_etl_execution_success() di akhir pipeline
- **End_ETL_Log_Failed_Step**: Transform step Execute SQL yang memanggil end_etl_execution_failed() dalam error handler
- **Execution_ID**: UUID unik yang dihasilkan oleh start_etl_execution() untuk tracking satu eksekusi ETL
- **LAST_ETL_RUN**: Variable timestamp yang digunakan untuk incremental loading filter
- **Etl_Log_Table**: Tabel audit trail di Warehouse yang menyimpan metadata eksekusi ETL
- **Row_Count**: Jumlah rows yang di-insert, update, atau delete dalam satu eksekusi ETL
- **Error_Handler**: Apache Hop error handling mechanism yang menangkap failure pada transform

## Requirements

### Requirement 1: Start ETL Execution Logging

**User Story:** As a Data Engineer, I want every pipeline execution to be logged with start timestamp and parameters, so that I can track when ETL runs and with what configuration

#### Acceptance Criteria

1. WHEN an ETL_Pipeline starts execution, THE Start_ETL_Log_Step SHALL call start_etl_execution() with etl_name 'MasterToWarehouse', pipeline_name matching the file name, and LAST_ETL_RUN parameter
2. THE Start_ETL_Log_Step SHALL capture the returned Execution_ID value into a pipeline variable
3. THE Start_ETL_Log_Step SHALL be named '00_Start_ETL_Log' and positioned as the first transform in the pipeline
4. THE Start_ETL_Log_Step SHALL use connection 'Warehouse' to execute the SQL function
5. THE Start_ETL_Log_Step SHALL record parameters as JSONB including LAST_ETL_RUN timestamp

### Requirement 2: End ETL Execution Success Logging

**User Story:** As a Data Engineer, I want successful pipeline executions to be logged with completion timestamp and row counts, so that I can verify data processing metrics

#### Acceptance Criteria

1. WHEN an ETL_Pipeline completes successfully, THE End_ETL_Log_Success_Step SHALL call end_etl_execution_success() with the captured Execution_ID
2. THE End_ETL_Log_Success_Step SHALL pass row count metrics (rows_inserted, rows_updated, rows_deleted) to the logging function
3. THE End_ETL_Log_Success_Step SHALL be named '99_End_ETL_Log_Success' and positioned as the final transform in the pipeline
4. THE End_ETL_Log_Success_Step SHALL use connection 'Warehouse' to execute the SQL function
5. THE End_ETL_Log_Success_Step SHALL update etl_metadata.last_successful_run automatically via the function

### Requirement 3: End ETL Execution Failed Logging

**User Story:** As a Data Engineer, I want failed pipeline executions to be logged with error details, so that I can troubleshoot failures quickly

#### Acceptance Criteria

1. WHEN an ETL_Pipeline encounters an error, THE End_ETL_Log_Failed_Step SHALL call end_etl_execution_failed() with the captured Execution_ID and error message
2. THE End_ETL_Log_Failed_Step SHALL be attached to error handling paths of critical transforms
3. THE End_ETL_Log_Failed_Step SHALL be named '99_End_ETL_Log_Failed'
4. THE End_ETL_Log_Failed_Step SHALL use connection 'Warehouse' to execute the SQL function
5. IF error message exceeds 500 characters, THEN THE End_ETL_Log_Failed_Step SHALL truncate the message to 500 characters

### Requirement 4: Dimension Pipeline Logging Integration

**User Story:** As a Data Engineer, I want all 7 dimension pipelines to have consistent logging, so that I can monitor dimension loading uniformly

#### Acceptance Criteria

1. THE ETL_Pipeline SHALL add logging steps to dim_pasien.hpl with pipeline_name 'dim_pasien'
2. THE ETL_Pipeline SHALL add logging steps to dim_anak.hpl with pipeline_name 'dim_anak'
3. THE ETL_Pipeline SHALL add logging steps to dim_ibu_hamil.hpl with pipeline_name 'dim_ibu_hamil'
4. THE ETL_Pipeline SHALL add logging steps to dim_lokasi.hpl with pipeline_name 'dim_lokasi'
5. THE ETL_Pipeline SHALL add logging steps to dim_posyandu.hpl with pipeline_name 'dim_posyandu'
6. THE ETL_Pipeline SHALL add logging steps to dim_petugas.hpl with pipeline_name 'dim_petugas'
7. THE ETL_Pipeline SHALL add logging steps to dim_waktu.hpl with pipeline_name 'dim_waktu'
8. WHEN a dimension pipeline is modified, THE ETL_Pipeline SHALL preserve existing incremental loading logic with LAST_ETL_RUN variable

### Requirement 5: Fact Pipeline Logging Integration

**User Story:** As a Data Engineer, I want both fact pipelines to have consistent logging, so that I can monitor fact loading uniformly

#### Acceptance Criteria

1. THE ETL_Pipeline SHALL add logging steps to fact_imunisasi.hpl with pipeline_name 'fact_imunisasi'
2. THE ETL_Pipeline SHALL add logging steps to fact_pemeriksaan.hpl with pipeline_name 'fact_pemeriksaan'
3. WHEN a fact pipeline is modified, THE ETL_Pipeline SHALL preserve existing dimension lookup logic
4. WHEN a fact pipeline is modified, THE ETL_Pipeline SHALL preserve existing incremental loading logic with LAST_ETL_RUN variable

### Requirement 6: Row Count Tracking

**User Story:** As a Data Engineer, I want accurate row counts logged for each execution, so that I can validate data completeness

#### Acceptance Criteria

1. THE End_ETL_Log_Success_Step SHALL track rows_inserted count from the Execute UPSERT transform
2. THE End_ETL_Log_Success_Step SHALL track rows_updated count from the Execute UPSERT transform
3. THE End_ETL_Log_Success_Step SHALL set rows_deleted to 0 for all MasterToWarehouse pipelines
4. WHERE row counts cannot be automatically captured, THE End_ETL_Log_Success_Step SHALL accept hardcoded placeholder values with a comment indicating manual adjustment needed

### Requirement 7: Pipeline Hop Ordering Preservation

**User Story:** As a Data Engineer, I want logging steps integrated without breaking existing data flow, so that ETL logic remains correct

#### Acceptance Criteria

1. WHEN logging steps are added to an ETL_Pipeline, THE ETL_Pipeline SHALL maintain the original hop connections between data transforms
2. THE Start_ETL_Log_Step SHALL not have input hops from other transforms
3. THE End_ETL_Log_Success_Step SHALL not send output to other transforms
4. THE ETL_Pipeline SHALL execute Start_ETL_Log_Step before any data extraction transform
5. THE ETL_Pipeline SHALL execute End_ETL_Log_Success_Step after all data loading transforms complete

### Requirement 8: Error Handler Configuration

**User Story:** As a Data Engineer, I want error handlers properly configured on critical transforms, so that failures are logged before pipeline termination

#### Acceptance Criteria

1. THE Execute UPSERT transform SHALL have error handling enabled with target transform '99_End_ETL_Log_Failed'
2. WHEN the Execute UPSERT transform fails, THE Error_Handler SHALL route the error to End_ETL_Log_Failed_Step
3. THE End_ETL_Log_Failed_Step SHALL capture error description from the error handling stream
4. WHEN End_ETL_Log_Failed_Step completes, THE ETL_Pipeline SHALL terminate with failure status

### Requirement 9: Execution ID Variable Management

**User Story:** As a Data Engineer, I want Execution_ID properly captured and propagated, so that all logging steps reference the same execution

#### Acceptance Criteria

1. THE Start_ETL_Log_Step SHALL define output field 'execution_id' with type String
2. THE Start_ETL_Log_Step SHALL use "Set Variables" or equivalent mechanism to store execution_id as pipeline variable
3. THE End_ETL_Log_Success_Step SHALL reference ${execution_id} variable in its SQL call
4. THE End_ETL_Log_Failed_Step SHALL reference ${execution_id} variable in its SQL call
5. WHEN execution_id variable is not available, THE End_ETL_Log steps SHALL log an error indicating variable capture failure

### Requirement 10: Database Connection Configuration

**User Story:** As a Data Engineer, I want all logging steps to use the correct database connection, so that logs are written to the Warehouse database

#### Acceptance Criteria

1. THE Start_ETL_Log_Step SHALL use database connection named 'Warehouse'
2. THE End_ETL_Log_Success_Step SHALL use database connection named 'Warehouse'
3. THE End_ETL_Log_Failed_Step SHALL use database connection named 'Warehouse'
4. WHEN connection 'Warehouse' is unavailable, THE logging step SHALL fail with connection error message
5. THE logging steps SHALL not modify or use the 'Master' database connection

### Requirement 11: SQL Function Parameter Formatting

**User Story:** As a Data Engineer, I want SQL function calls properly formatted with correct parameter types, so that PostgreSQL executes them without errors

#### Acceptance Criteria

1. THE Start_ETL_Log_Step SHALL format parameters JSON with proper JSONB casting using '::JSONB'
2. THE End_ETL_Log_Success_Step SHALL pass integer row counts without quotes
3. THE End_ETL_Log_Failed_Step SHALL escape single quotes in error messages using PostgreSQL escaping rules
4. THE logging steps SHALL use dollar-quoted strings or proper escaping for string parameters
5. WHEN LAST_ETL_RUN is empty string, THE Start_ETL_Log_Step SHALL convert it to '1900-01-01' using COALESCE or equivalent

### Requirement 12: Pipeline XML Structure Preservation

**User Story:** As a Data Engineer, I want modified pipeline files to remain valid Apache Hop XML, so that Hop GUI can open them without errors

#### Acceptance Criteria

1. WHEN an ETL_Pipeline file is modified, THE ETL_Pipeline SHALL maintain valid XML structure with proper encoding declaration
2. THE logging transforms SHALL use correct Apache Hop transform type identifiers (TableInput, ExecSql, ScriptValueMod)
3. THE logging transforms SHALL include required XML elements (name, type, description, distribute, copies, partitioning)
4. THE hop elements SHALL reference valid transform names in "from" and "to" attributes
5. WHEN a pipeline is saved, THE ETL_Pipeline SHALL validate XML schema compliance

### Requirement 13: Audit Trail Completeness

**User Story:** As a Compliance Officer, I want complete audit trail for all ETL executions, so that I can demonstrate data lineage compliance

#### Acceptance Criteria

1. THE Etl_Log_Table SHALL contain one record per ETL_Pipeline execution with status RUNNING, SUCCESS, or FAILED
2. THE Etl_Log_Table SHALL record start_time for every execution attempt
3. WHEN an execution succeeds, THE Etl_Log_Table SHALL record end_time, duration_seconds, and rows_total
4. WHEN an execution fails, THE Etl_Log_Table SHALL record end_time, duration_seconds, error_message, and status='FAILED'
5. THE Etl_Log_Table SHALL preserve records for at least 90 days per data retention policy

### Requirement 14: Incremental Loading Compatibility

**User Story:** As a Data Engineer, I want logging to work seamlessly with incremental loading, so that LAST_ETL_RUN metadata remains accurate

#### Acceptance Criteria

1. WHEN an ETL_Pipeline uses LAST_ETL_RUN variable in WHERE clause, THE Start_ETL_Log_Step SHALL log the LAST_ETL_RUN value to parameters field
2. WHEN end_etl_execution_success() completes, THE function SHALL update etl_metadata.last_successful_run to current timestamp
3. WHEN an ETL_Pipeline fails, THE function SHALL not update etl_metadata.last_successful_run
4. THE next execution SHALL use the last_successful_run from the previous successful execution
5. WHERE LAST_ETL_RUN is null or empty, THE ETL_Pipeline SHALL default to '1900-01-01' for full reload

### Requirement 15: Performance Impact Minimization

**User Story:** As a Data Engineer, I want logging overhead to be minimal, so that ETL execution time does not significantly increase

#### Acceptance Criteria

1. THE Start_ETL_Log_Step execution SHALL complete within 100 milliseconds
2. THE End_ETL_Log_Success_Step execution SHALL complete within 100 milliseconds
3. THE End_ETL_Log_Failed_Step execution SHALL complete within 100 milliseconds
4. THE logging steps SHALL execute database functions using prepared statement or equivalent optimization
5. WHEN total pipeline execution time is measured, THE logging steps SHALL contribute less than 5% overhead

### Requirement 16: Error Message Informativeness

**User Story:** As a Data Engineer, I want error messages to contain sufficient detail, so that I can diagnose failures without additional investigation

#### Acceptance Criteria

1. WHEN an Execute UPSERT transform fails, THE End_ETL_Log_Failed_Step SHALL capture the SQL error message
2. WHEN a dimension lookup fails, THE End_ETL_Log_Failed_Step SHALL capture the lookup key values
3. THE error_message field SHALL include transform name where the error occurred
4. THE error_message field SHALL include row number or identifier if available
5. WHERE multiple errors occur, THE End_ETL_Log_Failed_Step SHALL log the first error encountered

