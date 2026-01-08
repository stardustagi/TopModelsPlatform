package node_report

import (
	"github.com/stardustagi/TopLib/libs/server"
	"github.com/stardustagi/TopModelsPlatform/backend"
)

func (rep *NodeHttpReportService) initialization() {
	rep.app.AddGroup("report", server.Request(), backend.PlatformUserAccess())

	rep.app.AddPostHandler("report", server.NewHandler(
		"getDayReportList",
		[]string{"Report"},
		rep.GetDayReportList,
	))

	rep.app.AddPostHandler("report", server.NewHandler(
		"getDayReportSummary",
		[]string{"Report"},
		rep.GetDayReportSummary,
	))

	rep.app.AddPostHandler("report", server.NewHandler(
		"getReportDetailList",
		[]string{"Report"},
		rep.GetReportDetailList,
	))
}
