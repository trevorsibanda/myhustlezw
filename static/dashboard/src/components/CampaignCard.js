import React, { Component } from "react"
import { Link } from "react-router-dom"

class CampaignCard extends Component {
    render() {
        return (
            <Link className="card" to={"/creator/content/"+ this.props.campaign._id } style={{"backgroundColor": "ghostwhite"}} >
                <img className="card-img-top" src="/assets/img/blog/01.jpg" alt="Card image cap" />

                <div className="card-body">
                    <h4 className="card-title">{this.props.campaign.title}</h4>
                    <p className="card-text">{this.props.campaign.description.substr(0, 144)}{this.props.campaign.description.length > 144 ? '...':''}</p>
                </div>
        </Link>
        )
    }
}

export default CampaignCard;