import React, { Component } from "react"
import { NavLink } from 'react-router-dom'

class SideNav extends Component {

    render() {
        return (
            <ul id="sidemenu" className="nav nav-pills flex-column d-none d-md-block" >
                {this.props.user.verified ?
                    <><li className=" nav-king nav-item">
                        <NavLink exact={true} className="nav-link" to="/creator/dashboard"><i className="fa fa-home"></i> My Account</NavLink>
                    </li>
                    <li className="nav-item nav-king">
                            <NavLink className="nav-link" to="/creator/supporters"><i className="fa fa-heart"></i> My {this.props.user.page.supporter }s</NavLink>
                    </li></>
                    :
                    <><li className=" nav-king nav-item">
                        <NavLink exact={true} className="nav-link" to="/creator/dashboard/timeline"><i className="fa fa-newspaper"></i> Content Feed</NavLink>
                    </li>
                    <li className="nav-item nav-king">
                            <NavLink className="nav-link" to="/creator/supporters/subscriptions"><i className="fa fa-heart"></i> My Subscriptions</NavLink>
                    </li></>
                }
                    
                
                <li className="nav-item nav-king ">
                        <NavLink className="nav-link" to="/creator/content"><i className="fa fa-cube"></i> My Content</NavLink>
                </li>
                
                <li className="nav-item nav-king">
                    <NavLink className="nav-link" to="/creator/wallet"><i className="fa fa-credit-card"></i> My Wallet</NavLink>
                </li>
                <li className="nav-item nav-king">
                    <NavLink className="nav-link" to="/creator/settings"><i className="fa fa-cog"></i> Settings</NavLink>
                </li>
                <li className="nav-item nav-king">
                        <NavLink className="nav-link" to="/creator/help"><i className="fa fa-question"></i> Help and support</NavLink>
                </li>
                {this.props.user.verified ? <></> :
                    <li className="nav-item nav-king">
                        <NavLink className="nav-link bg-white" to="/creator/verify-identity">
                            <i className="fa fa-user text-success"></i> Verify your identity
                        </NavLink>
                    </li>
                }
                <li className="nav-item nav-king">
                    <a className="nav-link" href="/logout">Logout</a>
                </li>
            </ul>
        );
    }

}


export default SideNav;