import React, { Component } from "react"
import { NavLink, Switch, Route } from 'react-router-dom'
import Restrict from "./Restrict"




class BottomNav extends Component {
    constructor(props) {
        super(props)

        this.closeMenu = () => {
            document.getElementById('mobileNavBtn').click()
        }
    }

    render() {
        return (
            <nav class="navbar fixed-bottom navbar-expand-md custom-navbar navbar-container d-xs-block d-sm-block d-md-none">
                <div class="container">
                    <div class="mobile-app-icon-bar bg-dark" style={{'display': 'flex', 'width': '100%'}}>
                        {this.props.user.verified ?
                            <NavLink activeClassName="" activeStyle={{ "color": "burlywood" }} to="/creator/dashboard" strict className="btn" ><i class="fa fa-2x fa-home" aria-hidden="true"></i></NavLink> :
                            <NavLink activeClassName="" activeStyle={{ "color": "burlywood" }} to="/creator/dashboard/timeline" className="btn" ><i class="fa fa-2x fa-newspaper" aria-hidden="true"></i></NavLink>
                        }
                        <NavLink activeClassName="" activeStyle={{ "color": "burlywood" }} to="/creator/supporters" className="btn" ><i class="fa fa-2x fa-heart" aria-hidden="true"></i></NavLink>
                        <NavLink activeClassName="" activeStyle={{ "color": "burlywood" }} to="/creator/content" className="btn"  ><i class="fa fa-2x fa-plus-circle" aria-hidden="true"></i></NavLink>
                        <NavLink activeClassName="" activeStyle={{ "color": "burlywood" }} to="/creator/wallet" className="btn"  ><i class="fa fa-2x fa-credit-card" aria-hidden="true"></i></NavLink>
                        <NavLink activeClassName="" activeStyle={{ "color": "burlywood" }} to="/creator/settings" className="btn"  ><i class="fa fa-2x fa-cog" aria-hidden="true"></i></NavLink>
                    </div>
                </div>
            </nav>
        );
    }
}

export default BottomNav;

