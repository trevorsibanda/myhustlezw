import React, { Component } from "react";
import { Link, Redirect } from "react-router-dom"
import v1 from "../api/v1";


import SupportersList from "../components/ListSupporters"

class Supporters extends Component {

    constructor(props){
        super(props)
    }

    render() {
        return <>
            {this.props.user.verified ? 
            <div className="padding-top-50">
                <SupportersList title={"All " + this.props.user.page.supporter+ "s"} supporterName={this.props.user.page.supporter} type="recent" count={30} />
            </div>
                : <Redirect to="/creator/supporters/subscriptions" />}
            
            </>
    }
}

export default Supporters;