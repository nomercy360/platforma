import {Router} from "@solidjs/router";
import {FileRoutes} from "@solidjs/start/router";
import {Suspense} from "solid-js";
import Nav from "~/components/Nav";
import "./app.css";
import Controls from "~/components/controls";

export default function App() {
    return (
        <Router
            root={props => (
                <div class='bg-background min-h-screen'>
                    <Controls/>
                    <div class='flex flex-row items-start pr-4'>
                        <Nav/>
                        <Suspense>{props.children}</Suspense>
                    </div>
                </div>
            )}
        >
            <FileRoutes/>
        </Router>
    );
}
