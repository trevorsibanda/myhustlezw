import React, { Component } from "react"
import { Link } from "react-router-dom"


import ReactTimeAgo from 'react-time-ago'
import v1 from "../api/v1"
import money from "./payments/Amount"


class SubscriptionListCard extends Component {
    render() {
        let pic = v1.assets.profPicURL(this.props.subscription.creator)
        let target = this.props.fans ? "fan": "creator"
        return (
            <div class="col-md-12 media media-single" to={"/"+ target +"/subscriptions/" + this.props.subscription.creator._id}>
                <span class="avatar avatar-lg bg-danger"><Link to={"/@" + this.props.subscription.creator.username}><img className="avatar" src={ pic } alt="..." /></Link></span>
                <div class="media-body">
                    <h5><Link to={"/@" + this.props.subscription.creator.username}>@{this.props.subscription.creator.username}</Link></h5>
                    <p>{this.props.subscription.creator.profile.description.substr(0, 144)}{this.props.subscription.creator.profile.description.length > 144 ? '...' : ''}</p>
                    <table class="table table-responsive">
                        <tbody>
                            <tr>
                                <td>Paid</td>
                                <td>{money.format(this.props.subscription.sub.currency, this.props.subscription.sub.amount) }</td>
                            </tr>
                            <tr>
                                <td>Subscribed</td>
                                <td><ReactTimeAgo date={this.props.subscription.sub.created_at} /></td>
                            </tr>
                            <tr>
                                <td>Expires</td>
                                <td><ReactTimeAgo date={this.props.subscription.sub.expires} /></td>
                            </tr>
                        </tbody>
                    </table>
                    
                </div>
            </div>
        )
    }
}

export default SubscriptionListCard;