document.body.dataset['scheme'] = window.localStorage.getItem('scheme') || '';
document.body.dataset['theme'] = window.localStorage.getItem('theme') || '';

document.addEventListener('DOMContentLoaded', function() {
  var preloadStylesheetCallback = function () {
    this.removeEventListener('load', preloadStylesheetCallback);
    this.rel = 'stylesheet';
  };

  var lazyStylesheets = document.querySelectorAll('link[rel="preload"]');
  lazyStylesheets.forEach(function(lazyStylesheet) {
    lazyStylesheet.addEventListener('load', preloadStylesheetCallback);
    lazyStylesheet.href = lazyStylesheet.dataset.href;
  });

  document.querySelectorAll('a.theme-switcher').forEach(a => a.addEventListener('click', (e) => {
    e.preventDefault();
    document.body.dataset['theme'] = e.target.dataset['theme'];
    window.localStorage.setItem('theme', e.target.dataset['theme']);
  }));
  document.querySelectorAll('a.scheme-switcher').forEach(a => a.addEventListener('click', (e) => {
    e.preventDefault();
    document.body.dataset['scheme'] = e.target.dataset['scheme'];
    window.localStorage.setItem('scheme', e.target.dataset['scheme']);
  }));
});
