import {Timer} from "./types/app"

interface TimerUI {
    duration: string;
    show: boolean;
    running: boolean;
    done: boolean;
    display: string;
}

export default (timer: Timer): TimerUI => ({
    duration: timer.duration,
    show: timer.show,
    running: timer.running,
    done: timer.done,
    display: timer.display,
});