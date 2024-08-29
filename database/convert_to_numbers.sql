CREATE OR REPLACE FUNCTION convert_text_to_integer(val text)
RETURNS INTEGER
AS $body$
    BEGIN
        RETURN coalesce(nullif(val, ''), '0')::integer;
    END;
$body$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION convert_text_to_real(val text)
RETURNS REAL
AS $body$
    BEGIN
        RETURN coalesce(nullif(val, ''), '0')::real;
    END;
$body$ LANGUAGE plpgsql;


ALTER TABLE books
    --ALTER COLUMN text_reviews_count TYPE INTEGER
    --USING convert_varchar_to_integer(text_reviews_count),

    ALTER COLUMN num_pages TYPE INTEGER
    USING convert_varchar_to_integer(num_pages),

    ALTER COLUMN publication_day TYPE INTEGER
    USING convert_varchar_to_integer(publication_day),

    ALTER COLUMN publication_month TYPE INTEGER 
    USING convert_varchar_to_integer(publication_month),

    ALTER COLUMN publication_year TYPE INTEGER 
    USING convert_varchar_to_integer(publication_year),

    ALTER COLUMN ratings_count TYPE INTEGER 
    USING convert_varchar_to_integer(ratings_count),

    ALTER COLUMN average_rating TYPE REAL 
    USING convert_varchar_to_real(average_rating);