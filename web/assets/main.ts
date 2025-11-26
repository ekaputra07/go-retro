import Alpine from "alpinejs";
import { App } from "./application";

window.Alpine = Alpine;
Alpine.store("app", App);
Alpine.start();