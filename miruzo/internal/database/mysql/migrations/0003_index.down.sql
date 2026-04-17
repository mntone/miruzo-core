-- Drop index for imagelist
DROP INDEX ix_stats_engaged ON stats;
DROP INDEX ix_stats_hall_of_fame ON stats;
DROP INDEX ix_stats_first_love ON stats;
DROP INDEX ix_stats_recently ON stats;
DROP INDEX ix_ingests_chronological ON ingests;
DROP INDEX ix_images_latest ON images;

-- Drop index for action
DROP INDEX ix_actions_love_canceled_lookup ON actions;
DROP INDEX ix_actions_love_lookup ON actions;
