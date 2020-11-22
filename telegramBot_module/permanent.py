MODULE_NAME = "doeparser"

WEEKDAYS = ["MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"]

DATABASE_BACKUP_1 = "db.old.sqlite3"
DATABASE_BACKUP_2 = "db.old.old.sqlite3"

SCHEDULE_NAME_CSV = "schedule.csv"
SCHEDULE_NAME = "schedule.xlsx"
SCHEDULE_BACKUP_1 = "schedule.old.xlsx"
SCHEDULE_BACKUP_2 = "schedule.old.old.xlsx"

SCHEDULE_DOWNLOAD_LINK = "https://cdn.filesend.jp/private/i7geQHiEj-S2jelNnAQC8GHi6n6t2SbNe5S8u9q2FXTUe-z_PF0J0jAmz1mzuOYk/schedule.csv"
SCHEDULE_MIN_SIZE_BYTES = 20#30 * 1024
SCHEDULE_LAST_COLUMN = 36
SCHEDULE_LAST_ROW = 134

MESSAGE_ERROR_NOTIFY = "Schedule parse error occurred. Please check manually."
MESSAGE_ERROR_PARSE_SYNTAX = "Error during schedule parse with"
MESSAGE_ERROR_UNKNOWN_GROUP = "Unknown group found in schedule"

ADMIN_NOTIFY_TIME = "20:00"
ADMIN_NOTIFY_TABLE_CHANGES = True
