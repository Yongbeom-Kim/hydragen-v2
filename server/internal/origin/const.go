package origin

import "os"

var BACKEND_URL_PREFIX = os.Getenv("SITE_ADDR") + "/v2/api/"
