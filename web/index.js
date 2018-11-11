import "bootstrap/dist/css/bootstrap.min.css";
import $ from "jquery";
import "bootstrap";

var $sidebar = 0;

$(document).ready(function() {
	$("#toggleSidebar").click(function() {
		if ($sidebar === 1) {
			$("#sidebar").hide();
			$("#toggleSidebar i").addClass("glyphicon-chevron-left");
			$("#toggleSidebar i").removeClass("glyphicon-chevron-right");
			$sidebar = 0;
		}
		else {
			$("#sidebar").show();
			$("#toggleSidebar i").addClass("glyphicon-chevron-right");
			$("#toggleSidebar i").removeClass("glyphicon-chevron-left");
			$sidebar = 1;
		}

        return false;
	});
});