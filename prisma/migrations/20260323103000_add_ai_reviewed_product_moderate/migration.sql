-- Add product moderation state for items already checked by AI
ALTER TYPE "ProductModerate" ADD VALUE IF NOT EXISTS 'AI_REVIEWED';
