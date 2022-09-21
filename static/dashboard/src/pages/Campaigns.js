import React, { Component } from "react";
import { Link, Redirect } from "react-router-dom"
import v1 from "../api/v1";

import CampaignList from "../components/CampaignList";


class Campaigns extends Component {

    constructor(props) {
        super(props)

        v1.page.set({ title: 'Dashboard / Your Content' })
        v1.page.track()

        this.state = {
            campaigns: []
        }

        v1.campaign.listAll({}, false).then(campaigns => {
            this.setState({campaigns})
        }).catch(_ => {
            v1.campaign.listAll({}, true).then(campaigns => {
                this.setState({ campaigns })
            })
        })
    }

    render() {
        
        return !this.props.user.verified ? <Redirect to="/creator/content/new" /> : (
            <>
                <div class="box">
                    <div class="box-header">
                        <div class="row">
                            <div class="col-md-8">
                                <h4 class="box-title">My Content</h4>
                                <p>
                                    <small>All your content and services are conveniently listed here.
                                        
                                    </small>
                                </p>
                            </div>
                            <div class="col-md-4">
                                <div class="box-controls pull-right">
                                    <select class="form-control" >
                                        <option value="ZWL">All content</option>
                                        <option value="ZWL"></option>
                                        <optgroup label="Files and media" >
                                            <option >Private content (for subscribers)</option>
                                            <option>Pay per view/download</option>
                                            <option >Free files and media</option>
                                        </optgroup>
                                        <option>Services only</option>
                                    </select>
                                &nbsp;
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <CampaignList campaigns={this.state.campaigns} creator={this.props.user} />
             
            </>
        )
    }

}

export default Campaigns;