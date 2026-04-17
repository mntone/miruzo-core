-- sqlc workaround:
-- This file is marked as Down migration so sqlc ignores unsupported
-- CREATE TRIGGER statements. Do not move or modify this marker.
-- +migrate Down
CREATE TRIGGER ai_ingests_id_guard
AFTER INSERT ON ingests
FOR EACH ROW
BEGIN
	IF NEW.id > 9007199254740991/*2^53 - 1*/ THEN
		SIGNAL SQLSTATE '45000'
			SET MESSAGE_TEXT = 'ingests.id exceeds max safe integer';
	END IF;
END;
