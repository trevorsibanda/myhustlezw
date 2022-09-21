import React, { Component } from "react";
import { Link } from "react-router-dom"
import v1 from "../api/v1";

import ServiceListCard from "./ServiceListCard"


class ServicesList extends Component {
    render() {
        return (
            <div class="box">
                <div class="box-header with-border">
                    <h4 class="box-title">Services you are offering</h4>
                </div>
                <div class="box-body row media-list media-list-hover">
                    {this.props.campaigns.map((campaign, index) => {
                        return <><ServiceListCard index={index} campaign={campaign} /></>
                    })}
                </div>
            </div>
        )
    }
}

export default ServicesList