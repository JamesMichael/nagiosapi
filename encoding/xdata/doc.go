/*

	Package xdata provides routines for parsing nagios external state data files. It works in a manner similar to
	encoding/json.

	This format is used by nagios to store the current state of the system (statusdata.dat), and
	data that needs to be persisted across service restarts (retention.dat).

	Complete structures are provided to extract data from the statusdata file, or you can pass in a custom structure
	with fewer fields to parse only the relevant data for your application.

	The canonical implementation of xdata can be found at:
	https://github.com/NagiosEnterprises/nagioscore/tree/master/xdata

*/
package xdata
