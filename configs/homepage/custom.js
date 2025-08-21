/* Custom JavaScript for Bolabaden */
;(function(){
  // Ensure external links open in new tabs
  document.addEventListener('click', (e) => {
    const a = e.target.closest('a');
    if (!a) return;
    const isExternal = a.hostname && a.hostname !== window.location.hostname
    if (isExternal) {
      a.setAttribute('target', '_blank')
      a.setAttribute('rel', 'noopener noreferrer')
    }
  })
})();
