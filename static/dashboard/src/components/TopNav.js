import React, { Component } from "react"
import { NavLink, Switch, Route, Link } from 'react-router-dom'
import store from "store2"
import AuthTopNav from "./AuthTopNav"
import Restrict from "./Restrict"


class TopNav extends Component {
    constructor(props) {
        super(props)

        this.closeMenu = () => {
            document.getElementById('mobileNavBtn').click()
        }
    }

    render() {
        return (
            <Switch>
                <Route path="/auth" >
                    <AuthTopNav />
                </Route>
                <Route path="/" >
                    <nav class="navbar fixed-top navbar-expand-md custom-navbar navbar-container">
                        <div class="container">
                            <img class="navbar-brand" src="/assets/img/logo.svg" id="logo_custom" width="96px" alt="logo" />
                            <button class="navbar-toggler navbar-toggler-right custom-toggler collapsed" id="mobileNavBtn" type="button" data-toggle="collapse" data-target="#collapsibleNavbar" aria-expanded="false" >
                                <span class="navbar-toggler-icon " ></span>
                            </button>
                            <div class="navbar-collapse collapse" id="collapsibleNavbar" >
                                {this.props.user.logged_in ?
                                    <ul class="navbar-nav ml-auto d-xs-block d-sm-block d-md-none">
                                        <li className=" nav-king nav-item">
                                            <NavLink exact={true} className="nav-link" onClick={this.closeMenu} to="/creator/dashboard"><i className="fa fa-home"></i> My dashboard</NavLink>
                                        </li>
                                        <Restrict user={this.props.user}>
                                            <li className=" nav-king nav-item">
                                                <NavLink className="nav-link" rel="noreferrer" to={'/@' + this.props.user.username}><i className="fa fa-link"></i> My Page</NavLink>
                                            </li>
                                        </Restrict>
                                    
                                        <li className="nav-item nav-king">
                                            <NavLink className="nav-link" onClick={this.closeMenu} to="/creator/content/new"><i className="fa fa-plus-circle"></i> Create New content/service</NavLink>
                                        </li>
                                        <li className="nav-item nav-king">
                                            <a  href="/logout" className="nav-link" onClick={window.logout}>Logout</a>
                                        </li>
                                    </ul> :
                                    <ul class="navbar-nav ml-auto d-xs-block d-sm-block d-md-none">
                                        <li className=" nav-king nav-item">
                                            <a href="/" className="nav-link" onClick={this.closeMenu} ><i className="fa fa-home"></i> Home </a>
                                        </li>
                                        <li className="nav-item nav-king">
                                            <Link to="/auth/login" className="nav-link" onClick={this.closeMenu}><i class="fa fa-user"></i> Login/Create a MyHustle account</Link>
                                        </li>
                                    </ul>
                                }
                            </div>
                        </div>
                    </nav>
                </Route>
            </Switch>
            
        );
    }
}

export default TopNav;