var $rows = $('#omx-table td');
$('#search').keyup(function() {
  var val = $.trim($(this).val()).replace(/ +/g, ' ').toLowerCase().split(' ');

  $rows.hide().filter(function() {
    var text = $(this).text().replace(/\s+/g, ' ').toLowerCase();
    var matchesSearch = true;
    $(val).each(function(index, value) {
      matchesSearch = (!matchesSearch) ? false : ~text.indexOf(value);
    });
    return matchesSearch;
  }).show();
});