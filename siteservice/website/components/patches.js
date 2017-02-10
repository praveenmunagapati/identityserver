function initializePolyfills() {
  if (!String.prototype.startsWith) {
      String.prototype.startsWith = function (haystack, needle) {
          return haystack.lastIndexOf(needle, 0) === 0;
      };
  }
  if (!String.prototype.includes) {
      String.prototype.includes = function (search, start) {
          'use strict';
          if (typeof start !== 'number') {
              start = 0;
          }

          if (start + search.length > this.length) {
              return false;
          } else {
              return this.indexOf(search, start) !== -1;
          }
      };
  }
  if (!Array.prototype.find) {
      Object.defineProperty(Array.prototype, 'find', {
          value: function (predicate) {
              'use strict';
              if (this == null) {
                  throw new TypeError('Array.prototype.find called on null or undefined');
              }
              if (typeof predicate !== 'function') {
                  throw new TypeError('predicate must be a function');
              }
              var list = Object(this);
              var length = list.length >>> 0;
              var thisArg = arguments[1];
              var value;

              for (var i = 0; i < length; i++) {
                  value = list[i];
                  if (predicate.call(thisArg, value, i, list)) {
                      return value;
                  }
              }
              return undefined;
          }
      });
  }
}
