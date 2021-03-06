package version

// VERSION is supplied with the git committish this is built from
var VERSION = ""

// IMPORTANT: These versions are overridden by version ldflags specifications VERSION_VARIABLES in the Makefile

// DdevVersion is the current version of ddev, by default the git committish (should be current git tag)
var DdevVersion = "v0.3.0-dev" // Note that this is overridden by make

// WebImg defines the default web image used for applications.
var WebImg = "drud/nginx-php-fpm7-local" // Note that this is overridden by make

// WebTag defines the default web image tag for drud dev
var WebTag = "v0.4.0" // Note that this is overridden by make

// DBImg defines the default db image used for applications.
var DBImg = "drud/mysql-docker-local-57" // Note that this is overridden by make

// DBTag defines the default db image tag for drud dev
var DBTag = "v0.3.0" // Note that this is overridden by make

// DBAImg defines the default phpmyadmin image tag used for applications.
var DBAImg = "drud/phpmyadmin"

// DBATag defines the default phpmyadmin image tag used for applications.
var DBATag = "v0.2.0"

// RouterImage defines the image used for the router.
var RouterImage = "drud/nginx-proxy" // Note that this is overridden by make

// RouterTag defines the tag used for the router.
var RouterTag = "v0.3.0" // Note that this is overridden by make

// COMMIT is the actual committish, supplied by make
var COMMIT = "COMMIT should be overridden"

// BUILDINFO is information with date and context, supplied by make
var BUILDINFO = "BUILDINFO should have new info"
