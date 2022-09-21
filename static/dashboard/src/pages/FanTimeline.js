import React, { Component } from "react";
import { Link } from "react-router-dom"


import v1 from "../api/v1";

class FanTimeline extends Component {
    constructor(props) {
        super(props)
    }

    render() {
        return (
            <>
                <h5>Overview</h5>

                <div className="offer-single-item padding-bottom-40">
                    <div className="icon">
                        {this.props.user.fullname.substr(0, 1)}
                    </div>
                    <div className="content">
                        <h4 className="title"><a href="">{this.props.user.fullname}</a></h4>
                        {this.props.user.fullname}
                        <br />
                    </div>
                </div>
                <div className="row justify-content-center">
                    <div class="col-md-12">
                        <h4>Quick actions</h4>
                    </div>
                    <div className="col-md-4">
                        <Link to="/fan/content" className="btn btn-default btn-block"><i className="fa fa-newspaper"></i> Content Feed</Link>
                    </div>
                    <div className="col-md-4">
                        <Link to="/fan/subscriptions" className="btn btn-default btn-block"><i className="fa fa-heart"></i> My subscriptions ( 0 )</Link>
                    </div>
                    <div className="col-md-4">
                        <Link to="/fan/payments" className="btn btn-default btn-block"><i className="fa fa-credit-card"></i> Payments History</Link>
                    </div>
                    <div className="col-md-4">
                        <Link to="/fan/settings" className="btn btn-default btn-block"><i className="fa fa-cog"></i> Settings</Link>
                    </div>
                    <div className="col-md-4">
                        <Link to="/help" className="btn btn-default btn-block"><i className="fa fa-cog"></i> Help and support</Link>
                    </div>
                    <div className="col-md-4">
                        <Link to="/fan/upgrade" className="btn btn-danger btn-block"><i className="fa fa-heart"></i> Upgrade to a creator account</Link>
                    </div>
                </div>

                <div className="row justify-content-center">

                </div>
            </>
        );
    }
}

export default FanTimeline;

