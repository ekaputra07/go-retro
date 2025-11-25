import Alpine from "alpinejs";
import app from "./app";

window.Alpine = Alpine;
Alpine.data('app', app);
Alpine.start();