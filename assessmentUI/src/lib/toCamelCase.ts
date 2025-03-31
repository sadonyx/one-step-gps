/**
 * Recursively convert `snake_case` keys to `camelCase` throughout nested objects and arrays.
 * @param {any} obj - An array or object of any type. *Note: this function is able to handle any input type.*
 * @returns
 */
function keysToCamelCase(obj: any): any {
  // Handle null or undefined
  if (obj === null || obj === undefined) {
    return obj;
  }

  // Handle arrays
  if (Array.isArray(obj)) {
    return obj.map(keysToCamelCase);
  }

  // Handle objects
  if (typeof obj === 'object') {
    return Object.keys(obj).reduce((acc, key) => {
      // Convert snake_case to camelCase
      const camelKey = stringToCamelCase(key);

      // Recursively convert nested objects/arrays
      acc[camelKey] = keysToCamelCase(obj[key]);

      return acc;
    }, {} as any);
  }

  // Return primitive values as-is
  return obj;
}

function stringToCamelCase(str: string) {
  return str.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase());
}

export { keysToCamelCase, stringToCamelCase };
