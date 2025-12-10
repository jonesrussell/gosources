-- Add support for global selectors in sources table
ALTER TABLE sources
    ADD COLUMN use_global_selectors BOOLEAN NOT NULL DEFAULT true,
    ALTER COLUMN selectors DROP NOT NULL;

-- Add index for use_global_selectors for faster queries
CREATE INDEX IF NOT EXISTS idx_sources_use_global_selectors ON sources(use_global_selectors);

-- Add comment to explain the column
COMMENT ON COLUMN sources.use_global_selectors IS 'When true, merges global default selectors with source-specific selectors. Source-specific selectors override defaults.';
COMMENT ON COLUMN sources.selectors IS 'Source-specific selectors. When use_global_selectors is true, these override global defaults. Can be NULL to use only global defaults.';
