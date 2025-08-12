-- Enhanced banned words table to support context-based checking
-- This is a future enhancement script that can be run when ready

-- Add context-related columns to existing banned words functionality
ALTER TABLE banned_words_pending ADD COLUMN context_pattern VARCHAR(500) NULL COMMENT 'Pattern for context-based checking';
ALTER TABLE banned_words_pending ADD COLUMN context_type ENUM('exact', 'regex', 'word_boundary', 'sentence') DEFAULT 'word_boundary' COMMENT 'Type of context matching';
ALTER TABLE banned_words_pending ADD COLUMN severity ENUM('low', 'medium', 'high', 'critical') DEFAULT 'medium' COMMENT 'Severity of the error';

-- If you have an approved banned_words table, add the same columns:
-- ALTER TABLE banned_words ADD COLUMN context_pattern VARCHAR(500) NULL COMMENT 'Pattern for context-based checking';
-- ALTER TABLE banned_words ADD COLUMN context_type ENUM('exact', 'regex', 'word_boundary', 'sentence') DEFAULT 'word_boundary' COMMENT 'Type of context matching';
-- ALTER TABLE banned_words ADD COLUMN severity ENUM('low', 'medium', 'high', 'critical') DEFAULT 'medium' COMMENT 'Severity of the error';

-- Example usage would be:
-- INSERT INTO banned_words_pending (word, reason, context_pattern, context_type, severity) 
-- VALUES ('huse', 'Should be "huset" in definite form', 'i huse', 'word_boundary', 'medium');

-- This allows for more sophisticated checking like:
-- - "huse" is only wrong when used as "i huse" (should be "i huset")
-- - "og" vs "å" confusion in infinitive contexts
-- - Preposition errors that depend on following words