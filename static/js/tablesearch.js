function searchTable() {
  
    var search_term = document.getElementById("search").value.toUpperCase();
    var tr = document.getElementById("omx-table").getElementsByTagName("tr");
  
    for (row of tr) {
        td1 = row.getElementsByTagName("td")[0].innerHTML.toUpperCase();
        td2 = row.getElementsByTagName("td")[1].innerHTML.toUpperCase();
        if (td1.indexOf(search_term) <= -1 && td2.indexOf(search_term) <= -1) {
            row.style.display = "none";
        } else {
            row.style.display = "";
        }
    }
}