import React, { Component } from "react";
import { Link } from "react-router-dom"
import v1 from "../api/v1";

import CampaignListCard from "./CampaignListCard"


class CampaignList extends Component {
    render() {
        return (
            <div class="box">
                <div class="box-body row media-list media-list-hover">
                    {this.props.campaigns.map((campaign, index) => {
                        return <><CampaignListCard index={index} creator={this.props.creator} campaign={campaign} /></>
                    })}
                </div>
            </div>
        )
    }
}

export default CampaignList