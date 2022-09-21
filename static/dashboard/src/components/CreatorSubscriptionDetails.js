import React, { Component } from "react"
import { Link } from "react-router-dom/cjs/react-router-dom.min";
import ReactTimeAgo from "react-time-ago/commonjs/ReactTimeAgo";

function CreatorSubscriptionDetails(props) {
    return (
        <div class="alert alert-primary">
            <div class="alert-body">
                <p>
                    <i class="fa fa-unlocked text-success"></i>
                    You subscribed to @{props.creator.username}'s account <ReactTimeAgo date={props.subscription.created_at} /> and your subscription expires
                    &nbsp; <ReactTimeAgo date={props.subscription.expires} />
                </p>
                <p>
                    <Link class="btn btn-block btn-black" to={"/creator/supporters/subscriptions"} >Your subscriptions</Link>
                </p>
            </div>
        </div>
    )
}

export default CreatorSubscriptionDetails;