import React, { Component } from "react";
import v1 from "../api/v1";
import money from "../components/payments/Amount"
import { Link } from "react-router-dom"


import ReactTimeAgo from 'react-time-ago'
import Preloader from "../components/PreLoader";

class FanViewSubscription extends Component {

    constructor(props) {
        super(props)

        this.state = {
            loaded: false,
            subscription: {}
        }

        let { id } = this.props.match.params


        v1.subscriptions.get(id).then(subscription => {
            this.setState({ subscription, loaded: true })
        }).catch(_ => {
            this.setState({loaded: false})
        })
    }

    render() {
        let target = this.props.fan ? "fan" : "creator"

        return this.state.loaded ?  (
            <>
                <div class="box">
                    <div class="box-header">
                        <div class="row">
                            <div class="col-md-8">
                                <h4 class="box-title">Your Subscription</h4>
                            </div>
                        </div>
                    </div>
                </div>
                <span class="avatar avatar-lg bg-danger"><img className="avatar" src={v1.assets.profPicURL(this.state.subscription.creator)} alt="..." /></span>
                <div class="media-body">
                    <h5>@{this.state.subscription.creator.username}</h5>
                    <p>is {this.state.subscription.creator.profile.description.substr(0, 144)}{this.state.subscription.creator.profile.description.length > 144 ? '...' : ''}</p>
                    <table class="table table-responsive">
                        <tbody>
                            <tr>
                                <td>Paid</td>
                                <td>{money.format(this.state.subscription.sub.currency, this.state.subscription.sub.amount)}</td>
                            </tr>
                            <tr>
                                <td>Subscribed</td>
                                <td><ReactTimeAgo date={this.state.subscription.sub.created_at} /></td>
                            </tr>
                            <tr>
                                <td>Expires</td>
                                <td><ReactTimeAgo date={this.state.subscription.sub.expires} /></td>
                            </tr>
                        </tbody>
                    </table>

                </div>

            </>
        ) : <Preloader />
    }

}

export default FanViewSubscription;