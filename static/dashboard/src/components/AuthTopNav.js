import { Component } from "react";

class AuthTopNav extends Component {
    render() {
        return (
            <header class="header">
                <nav class="navbar navbar-expand-lg navbar-light py-3">
                    <div class="container">
                        <a href="/" >
                            <img class="navbar-brand" src="/assets/img/logo.svg" id="logo_custom" alt="logo" width="150" />
                        </a>
                    </div>
                </nav>
            </header>
        )
    }
}

export default AuthTopNav;