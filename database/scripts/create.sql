CREATE TABLE IF NOT EXISTS tblVideos(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Title TEXT NOT NULL DEFAULT 'Unavailable',
	Description TEXT NOT NULL DEFAULT 'Unavailable',
	DurationSeconds INTEGER,
	OriginalURL TEXT NOT NULL,
	WebpageURL TEXT NOT NULL,
	LiveStatus TEXT,
	Availability TEXT,
    YoutubeViewCount INTEGER DEFAULT 0,
    LikeCount INTEGER DEFAULT 0,
    DislikeCount INTEGER DEFAULT 0,
    License TEXT,    
    AgeLimit INTEGER,
    PlayableInEmbed TEXT,
    UploadDate TEXT,
    ReleaseTimestamp INTEGER,
    ModifiedTimestamp INTEGER,
	IsFileDownloaded INTEGER NOT NULL DEFAULT 0,
	FileId INTEGER NOT NULL DEFAULT 0,
	ChannelId INTEGER,
	DomainId INTEGER,
	FormatId INTEGER,	
	YoutubeVideoId INTEGER,
	WatchCount INTEGER,
	IsDeleted INTEGER,
	CreatedDate INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS tblChannels(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Name TEXT NOT NULL DEFAULT 'Unavailable',
    ChannelFollowerCount INTEGER DEFAULT 0,
	ChannelURL TEXT,
	YoutubeChannelId INTEGER,
	CreatedDate INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS tblLiveStatusType(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	LiveStatus TEXT
);

CREATE TABLE IF NOT EXISTS tblAvailabilityType(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Availability TEXT
);

CREATE TABLE IF NOT EXISTS tblPlaylists(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Title TEXT NOT NULL DEFAULT 'Unavailable',
	ItemCount INTEGER,
    PlaylistChannel TEXT,
    PlaylistChannelId TEXT,
    PlaylistUploader TEXT,
    PlaylistUploaderId TEXT,
	ThumbnailFileId INTEGER,                    
	YoutubePlaylistId INTEGER,
	CreatedDate INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS tblPlaylistVideoFiles(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	VideoId INTEGER NOT NULL,
    PlaylistId INTEGER NOT NULL,
    PlaylistVideoIndex INTEGER NOT NULL,
	FileId INTEGER,
	ThumbnailFileId INTEGER,                    
	CreatedDate INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS tblFormats(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Format TEXT,
	FormatNote TEXT,
	Resolution TEXT,
	StreamType TEXT,
	CreatedDate INTEGER
);

CREATE TABLE IF NOT EXISTS tblDomains(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Domain TEXT DEFAULT 'Unavailable',
	CreatedDate INTEGER
);

CREATE TABLE IF NOT EXISTS tblTags(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Name TEXT,
	IsUsed INTEGER NOT NULL DEFAULT 0,	
	CreatedDate INTEGER
);

CREATE TABLE IF NOT EXISTS tblVideoFileTags(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	TagId INTEGER NOT NULL,
	VideoId INTEGER,
	FileId INTEGER,
	CreatedDate INTEGER
);

CREATE TABLE IF NOT EXISTS tblCategories(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Name TEXT,
	IsUsed INTEGER NOT NULL DEFAULT 0,	
	CreatedDate INTEGER
);

CREATE TABLE IF NOT EXISTS tblVideoFileCategories(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	CategoryId INTEGER NOT NULL,
	VideoId INTEGER,
	FileId INTEGER,
	CreatedDate INTEGER
);

CREATE TABLE IF NOT EXISTS tblFiles(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	VideoId INTEGER NOT NULL,
	FileType TEXT,
	Source TEXT,
	FilePath TEXT,
	FileName TEXT,
	Extension TEXT,
	FileSize INTEGER,
	FileSizeUnit TEXT,
	NetworkPath TEXT,
	IsDeleted INTEGER,
	CreatedDate INTEGER
);

CREATE TABLE IF NOT EXISTS tblSourceType(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	Source TEXT
);

CREATE TABLE IF NOT EXISTS tblFileType(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	File TEXT
);

-- logs all yt-dlp commands here
CREATE TABLE IF NOT EXISTS tblNetworkActivityLogs(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	ActivityTypeId INTEGER NOT NULL,	
	InputURL TEXT,
	Command TEXT,
	Result TEXT,
	CreatedDate INTEGER
);

-- logs all crud operations here
CREATE TABLE IF NOT EXISTS tblLocalActivityLogs(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	ActivityTypeId INTEGER NOT NULL,	
	VideoId INTEGER,
	FileId INTEGER,
	Remarks TEXT,
	CreatedDate INTEGER
);

CREATE TABLE IF NOT EXISTS tblActivityType(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	ActivityName TEXT,
	IsNetworkActivity INTEGER NOT NULL DEFAULT 0,
	IsActive INTEGER NOT NULL DEFAULT 1
);

-- It helps debug perf and errors.
-- It ties all activities performed(or not performed) and when.
CREATE TABLE IF NOT EXISTS tblAPILogs(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	APIName TEXT,
	ExecutionStart INTEGER,
	ExecutionEnd INTEGER,
	APIInputs TEXT,
	APIResult TEXT,
	ActivityLogId INTEGER,
	VideoId INTEGER,
	CreatedDate INTEGER
);

CREATE TABLE IF NOT EXISTS tblAPIType(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	APIName TEXT,
	APIVersion TEXT NOT NULL DEFAULT 'v1',
	IsActive INTEGER NOT NULL DEFAULT 1
);


-------------------------------------------------------------
-------------------------------------------------------------
--UPDATE-----------------------------------------------------

-- api will populate this
CREATE TABLE IF NOT EXISTS tblRequests(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	RequestUrl TEXT  NOT NULL, --Url or Identifier
	RequestType TEXT, --(Video, Channel, Playlist, Unknown)
	Metadata INTEGER, --(boolean)
	Thumbnail INTEGER, --(boolean)
	Content INTEGER, --(boolean)
	ContentFormat TEXT, --(Audio-Best, Video-Best, FormatId)
	Subtitles INTEGER, --(boolean)
	SubtitlesLanguage TEXT, --(Auto or provide language)
	IsProxied INTEGER, --(boolean)
	Proxy TEXT, -- Tor URL or CustomProxy		
    Scheduled INTEGER NOT NULL, --(Request created by scheduler, not to be executed by API)
	CreatedDate INTEGER,
	ModifiedDate INTEGER
);

CREATE TABLE IF NOT EXISTS tblRequestStatus(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
    RequestId INTEGER,
	RequestStatus INTEGER NOT NULL, -- Recieved, Analyzing, Queued, PartialComplete, Complete, Failed. 
	CreatedDate INTEGER,
	ModifiedDate INTEGER
);

-- Expand Requests and dump into queue
CREATE TABLE IF NOT EXISTS tblRequestQueue(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	RequestId INTEGER NOT NULL, 
	ContentId TEXT, --(YTVideoId)
	ProcessStatus TEXT, --(Queued, Started, Processing Metadata, Downloading, Done, Failed)
	RetryCount INTEGER, --(Manual retries only)
	Message TEXT, --(store failure reason or 'completed successfully')
	Cancelled INTEGER, --(Soft delete if you want to cancel queued download)    
    RequestType TEXT, -- (api-request/scheduled-request)
	CreatedDate INTEGER,
	ModifiedDate INTEGER
);


CREATE TABLE IF NOT EXISTS tblProcess(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	ProcessId INTEGER NOT NULL, --(pid for extractor, populated by API)
	RequestId INTEGER NOT NULL, --(Current RequestId that is being processed)
	Status INTEGER, --(boolean)[0 - Stopped, 1 - Processing]
	StartTime INTEGER,
	EndTime INTEGER,
	ProcessType TEXT, -- (api-request-process/scheduled-request-process)
	Message TEXT, --(If Panic is there, recover and log, if repeats, log and exit)
	CreatedDate INTEGER	
);

CREATE TABLE IF NOT EXISTS tblDownloadStatistics(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	RequestId INTEGER NOT NULL, 
	ContentId TEXT NOT NULL, 
	DownloadString TEXT,
	IsDownloadComplete INTEGER, --(boolean) [1 - yes, 2 - no]
	CreatedDate INTEGER,
	ModifiedDate INTEGER
);

-- Populated by API
-- Read by extractor
-- Disabled schedules will use the SoftDelete in RequestQueue
CREATE TABLE IF NOT EXISTS tblSchedule(
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	RequestId INTEGER, --(Scheduled Request)
	FriendlyName TEXT, --(friendly name or identifier -> default to guid)
	Type TEXT, --(Channel or Playlist)
	CRONExpression TEXT, --(default to update once every day at 2AM)	
	IsEnabled INTEGER, --(boolean) [1 - yes, 2 - no]
	IsDeleted INTEGER, --(boolean)[Deletes CRON Schedule]
	CreatedDate INTEGER
);

-- script will execute
-- execution will produce these logs which will then be shown on the UI
CREATE TABLE IF NOT EXISTS tblScheduleRunLogs (
	Id INTEGER PRIMARY KEY AUTOINCREMENT,
	ScheduleId INTEGER NOT NULL,
	StartTime INTEGER, 
	EndTime INTEGER,
	Status TEXT, --(Started, Running, Completed, Failed)
	Message TEXT, --(failure details or completed successfully) 
	CreatedDate INTEGER,
	ModifiedDate INTEGER
);