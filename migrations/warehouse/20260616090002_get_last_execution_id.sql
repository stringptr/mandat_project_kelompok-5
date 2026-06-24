CREATE OR REPLACE FUNCTION public.get_last_execution_id(p_etl_name VARCHAR) RETURNS VARCHAR AS $$
DECLARE v_execution_id VARCHAR(50);
BEGIN
    SELECT last_execution_id INTO v_execution_id FROM public.etl_metadata WHERE etl_name = p_etl_name;
    RETURN v_execution_id;
END;
$$ LANGUAGE plpgsql;
