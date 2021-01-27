export const apiUrl = 'http://localhost:10000';
export const socialWSUrl = 'ws://localhost:10000/social/ws?token=';
export const messengerWSUrl = 'ws://localhost:10000/messenger/ws?token=';
export const headers = {
  'Content-Type': 'application/json',
};

// Returns a function, that, as long as it continues to be invoked, will not
// be triggered. The function will be called after it stops being called for
// N milliseconds. If `immediate` is passed, trigger the function on the
// leading edge, instead of the trailing.
export function debounce(func, wait, immediate) {
  let timeout;
  return function() {
    const context = this, args = arguments;
    const later = function () {
      timeout = null;
      if (!immediate) func.apply(context, args);
    };
    const callNow = immediate && !timeout;
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
    if (callNow) func.apply(context, args);
  };
}
