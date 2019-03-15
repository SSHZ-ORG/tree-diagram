package paths

const (
	APIRenderEventPath = "/api/renderEvent"
	APIRenderPlacePath = "/api/renderPlace"
	APIRenderActorPath = "/api/renderActor"
	APIQueryEventsPath = "/api/queryEvents"

	CommandEnqueueDateRangePath   = "/admin/command/enqueueDateRange"
	CommandCompressEventSnapshots = "/admin/command/compressEventSnapshots"

	CrawlDatePath  = "/admin/crawl/date"
	CrawlActorPath = "/admin/crawl/actor"

	CronDailyPath   = "/admin/cron/daily"
	CronRevivePath  = "/admin/cron/revive"
	CronUndeadPath  = "/admin/cron/undead"
	CronCleanupPath = "/admin/cron/cleanup"
	CronExportPath  = "/admin/cron/export"
)
