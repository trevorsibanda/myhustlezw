import React, { Component } from "react"
import { Link } from "react-router-dom"


import ReactTimeAgo from 'react-time-ago'


class CampaignListCard extends Component {
    render() {
        let icon = 'fa-file'
        switch (this.props.campaign.type) {
            case 'video':
                icon = 'fa-film';
                break;
            case 'audio':
                icon = 'fa-music';
                break;
            case 'image':
                icon = 'fa-gallery';
                break;
            case 'service':
                icon = 'fa-vcard'
                break;
            case 'embed':
                icon = 'fa-youtube'
                break;
            default:
                break;
        }
        icon = 'fa '+icon
        return (
            <Link class="col-md-12 media media-single" to={"/@"+this.props.creator.username +"/"+this.props.campaign.uri}>
                <span class="avatar avatar-lg bg-danger"><i class={icon} style={{'color': 'white'}}></i></span>
                <div class="media-body">
                    <h5>{this.props.campaign.title}</h5>
                    <p>{this.props.campaign.description.substr(0, 144)}{this.props.campaign.description.length > 144 ? '...' : ''}</p>
                    <span><ReactTimeAgo date={this.props.campaign.created_at} /></span>
                </div>
            </Link>
        )
    }
}

export default CampaignListCard;