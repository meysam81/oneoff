// Simple cache utility with TTL support
const CACHE_PREFIX = "oneoff_cache_";
const DEFAULT_TTL = 5 * 60 * 1000; // 5 minutes

export class Cache {
  static set(key, value, ttl = DEFAULT_TTL) {
    const item = {
      value,
      expiry: Date.now() + ttl,
    };
    try {
      localStorage.setItem(CACHE_PREFIX + key, JSON.stringify(item));
    } catch (e) {
      console.warn("Cache write failed:", e);
    }
  }

  static get(key) {
    try {
      const itemStr = localStorage.getItem(CACHE_PREFIX + key);
      if (!itemStr) return null;

      const item = JSON.parse(itemStr);
      if (Date.now() > item.expiry) {
        localStorage.removeItem(CACHE_PREFIX + key);
        return null;
      }

      return item.value;
    } catch (e) {
      console.warn("Cache read failed:", e);
      return null;
    }
  }

  static remove(key) {
    try {
      localStorage.removeItem(CACHE_PREFIX + key);
    } catch (e) {
      console.warn("Cache remove failed:", e);
    }
  }

  static clear() {
    try {
      Object.keys(localStorage)
        .filter((key) => key.startsWith(CACHE_PREFIX))
        .forEach((key) => localStorage.removeItem(key));
    } catch (e) {
      console.warn("Cache clear failed:", e);
    }
  }
}

// Request deduplication
const pendingRequests = new Map();

export async function dedupe(key, fn) {
  if (pendingRequests.has(key)) {
    return pendingRequests.get(key);
  }

  const promise = fn()
    .then((result) => {
      pendingRequests.delete(key);
      return result;
    })
    .catch((error) => {
      pendingRequests.delete(key);
      throw error;
    });

  pendingRequests.set(key, promise);
  return promise;
}
