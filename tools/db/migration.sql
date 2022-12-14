DROP TABLE IF EXISTS user_tab;
DROP TABLE IF EXISTS map_tab;
DROP TABLE IF EXISTS game_tab;
DROP TABLE IF EXISTS game_user_tab;

CREATE TABLE user_tab(
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(64) NOT NULL,
    rating INT NOT NULL DEFAULT 0,
    moves_made BIGINT UNSIGNED NOT NULL DEFAULT 0,
    games_played INT UNSIGNED NOT NULL DEFAULT 0,
    highest_campaign INT NOT NULL DEFAULT 0,
    curr_campaign_game_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    time_created BIGINT NOT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX unq_email (email),
    UNIQUE INDEX unq_username (username)
);

CREATE TABLE map_tab(
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    type TINYINT UNSIGNED NOT NULL,
    height TINYINT UNSIGNED NOT NULL,
    width TINYINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    player_count TINYINT UNSIGNED NOT NULL,
    terrain_info BLOB,
    unit_info BLOB,
    author_user_id BIGINT UNSIGNED NOT NULL,
    is_campaign BOOLEAN NOT NULL DEFAULT FALSE,
    stat_play_count INT UNSIGNED NOT NULL DEFAULT 0,
    time_created BIGINT NOT NULL,
    time_modified BIGINT NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_author_user_id (author_user_id),
    INDEX idx_is_campaign (is_campaign, id)
);

CREATE TABLE game_tab(
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    type TINYINT UNSIGNED NOT NULL,
    height TINYINT UNSIGNED NOT NULL,
    width TINYINT UNSIGNED NOT NULL,
    player_count TINYINT UNSIGNED NOT NULL,
    terrain_info BLOB,
    unit_info BLOB,
    map_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(64) NOT NULL,
    creator_user_id BIGINT UNSIGNED NOT NULL,
    status TINYINT NOT NULL DEFAULT 0,
    turn_count INT NOT NULL DEFAULT 1,
    turn_player TINYINT NOT NULL DEFAULT 0,
    time_created BIGINT NOT NULL,
    time_modified BIGINT NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_status_created (status, time_created)
);

CREATE TABLE game_user_tab(
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    game_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    player_order TINYINT UNSIGNED NOT NULL,
    final_rank TINYINT UNSIGNED NOT NULL DEFAULT 0,
    final_turns INT NOT NULL DEFAULT 0,
    moves_made INT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    INDEX idx_game_id (game_id, player_order),
    INDEX idx_user_game_id (user_id, game_id)
);
