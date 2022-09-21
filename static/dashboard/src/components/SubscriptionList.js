import React, { Component } from "react";

import SubscriptionListCard from "./SubscriptionListCard"


class SubscriptionList extends Component {
    render() {
        return (
            <div class="box">
                <div class="box-body row media-list media-list-hover">
                    {this.props.subscriptions.map((subscription, index) => {
                        return <SubscriptionListCard index={index} fans={this.props.fans} subscription={subscription} />
                    })}
                </div>
            </div>
        )
    }
}

export default SubscriptionList