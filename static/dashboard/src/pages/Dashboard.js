import React, { Component } from "react";
import {Link, Redirect} from "react-router-dom"


import v1 from "../api/v1";

import SupportersList from "../components/ListSupporters"
import Preloader from "../components/PreLoader";
import Wallet from "./Wallet";

class Dashboard extends Component {
    constructor(props) {
        super(props)


        this.state = {
            ready: false,
        }

        v1.user.current(false).then(user => {
            this.setState({user, ready: true})
        }).catch(_ => {
            return v1.user.current(true).then(user => {
                this.setState({ user, ready: true })
            })
        })

    }

    render() {
        return !this.state.ready ? <Preloader /> : (this.state.user.verified ?  (
            <>
                 
                <h5>Overview</h5>
                
                <div className="offer-single-item padding-bottom-40">
                    <div className="icon">
                        {this.state.user.fullname.substr(0, 1)}
                    </div>
                    <div className="content">
                        <h4 className="title"><a href="">@{this.state.user.username}</a></h4>
                        {this.state.user.fullname}
                        <br />
                        {this.state.user.subscriptions.count} subscribers
                        <br />
                        <a href={ 'https://myhustle.co.zw/@'+this.state.user.username } target="_blank">https://myhustle.co.zw/@{this.state.user.username}</a>
                    </div>
                    <div className="action-buttons">
                        <Link to="/creator/page" className="btn btn-block btn-sm btn-success"><i className="fa fa-edit"></i> Edit page</Link>
                        <button className="btn btn-block btn-sm btn-success"><i className="fa fa-clipboard"></i> Copy link</button>
                    </div>
                </div>
                <h5>My links</h5>
                <div class="box">
                    <div class="box-header with-border">
                        <h4 class="box-title">My links</h4>
                    </div>
                    <div class="box-body p-0">
                        <div class="media-list media-list-hover media-list-divided">
                            <Link to='/auth/login' class="media media-single">
                                <i class="font-size-18 mr-0 flag-icon flag-icon-us"></i>
                                <span class="title">My page </span>
                                <span class="badge badge-pill badge-secondary">https://myhustle.co.zw/@{this.state.user.username}</span>
                            </Link>

                            <a class="media media-single" href="#">
                                <i class="font-size-18 mr-0 flag-icon flag-icon-ba"></i>
                                <span class="title">Subscribe</span>
                                <span class="badge badge-pill badge-primary">https://myhustle.co.zw/@{this.state.user.username}/subscribe</span>
                            </a>

                            <a class="media media-single" href="#">
                                <i class="font-size-18 mr-0 flag-icon flag-icon-ch"></i>
                                <span class="title">Buy me a {this.state.user.page.donation_item}</span>
                                <span class="badge badge-pill badge-info">https://myhustle.co.zw/@{this.state.user.username}/buymeacoffee</span>
                            </a>

                            
                        </div>
                    </div>
                </div>
                <div className="row justify-content-center">
                    <div className="col-md-4">
                        <Link to="/creator/content/new" className="btn btn-danger btn-block"><i className="fa fa-heart"></i> Create a new campaign</Link>
                    </div>
                </div>
                <Wallet user={this.state.user} hideRecent hideWithdraw hidePayout />
                <SupportersList title="Recent supporters"  type="recent" supporterName={this.state.user.page.supporter} count={10} />
            </>
        ) : <Redirect to="/creator/dashboard/timeline" />);
    }
}

export default Dashboard;

