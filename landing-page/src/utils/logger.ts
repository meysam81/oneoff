import log from "loglevel";

log.setDefaultLevel(log.levels.INFO);

if (import.meta.env.DEV) {
  log.setLevel(log.levels.DEBUG);
}

if (import.meta.env.PUBLIC_LOG_LEVEL) {
  const level = import.meta.env.PUBLIC_LOG_LEVEL.toUpperCase();
  if (log.levels[level] !== undefined) {
    log.setLevel(log.levels[level]);
  } else {
    log.warn(`Invalid log level: ${level}`);
  }
}

export default log;
