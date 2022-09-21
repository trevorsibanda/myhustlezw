import React, { Component } from "react";


import RichEditor from "../components/RichEditor"
import LiteYouTubeEmbed from "react-lite-youtube-embed";

import v1 from "../api/v1";
import { Redirect } from "react-router-dom";

class CreateYoutubeEmbedCampaign extends Component {
    constructor(props) {
        super(props)

        this.state = {
            multipleFiles: false,
            navigate: false,
            youtube_id: '',
            youtube_url: '',
            valid_yt_link: false,
            campaignimgurl: "https://img2.goodfon.com/wallpaper/nbig/d/31/art-gorod-budushie-razvaleny.jpg",
            campaign: {
                type: 'embed',
                title: '',
                price: 0.00,
                description: '  ',
                content: [],
                subscription: 'public',
                expires: 'never'
            }
        }

        this.createNewCampaign = this.createNewCampaign.bind(this)

        this.getYoutubeID = (link) => {
            return link.match(/(?:youtu\.be\/|youtube\.com(?:\/embed\/|\/v\/|\/watch\?v=|\/user\/\S+|\/ytscreeningroom\?v=|\/sandalsResorts#\w\/\w\/.*\/))([^\/&]{10,12})/)[1];
        }

        this.applyYoutubeLink = (link) => {
            try{
                let id = this.getYoutubeID(link)
                if (id.length > 0) {
                    this.setState({youtube_id: id, valid_yt_link: true, youtube_url: link})
                    let c = this.state.campaign
                    c.content = [id]
                    this.setState({campaign: c})
                }
            }catch(e){
                this.setState({valid_yt_link: false, youtube_url: link})
                console.log(e)
            }
        }
    }


    createNewCampaign() {

                if(this.state.youtube_id === ''){
                    return alert('You must specify a Youtube video.')
                }
        
                if(this.state.campaign.title.length < 10 ){
                    return alert('Title should be at least 10characters long')
                }
        v1.campaign.createCampaign(this.state.campaign, 'embed').then(c => {
            if (c._id) {
                //notify of creation success
                this.setState({navigate: true, campaign: c})
            } else {
                alert('Failed to create new the youtube embed\n\nErrors:'+c.errors)
            }
        })
        
    }
    render() {
        return this.state.navigate ? <Redirect to={'/creator/content/'+this.state.campaign._id} /> : (

            <div className="row">
                <div class="col-12">
                    <h4>Embed a youtube video.</h4>
                    <p><small>Already have content on Youtube? Just embed it!</small></p>
                </div>
                <div className="col-md-6" >
                    <LiteYouTubeEmbed
                        id={this.state.youtube_id} 
                        adNetwork={false} // Default true, to preconnect or not to doubleclick addresses called by YouTube iframe (the adnetwork from Google)
                        poster="hqdefault" 
                        title="YouTube Embedded video" 
                        noCookie={true} //Default false, connect to YouTube via the Privacy-Enhanced Mode using https://www.youtube-nocookie.com
                    />
                    <p>
                        <div className="form-group">
                            <label>Youtube Video Link</label>
                            <div className="input-group mb-3">
                                <div className="input-group-prepend">
                                    <span className="input-group-text"><i class="fa fa-youtube text-google"></i></span>
                                </div>
                                <input type="text" className="form-control" placeholder="https://www.youtube.com/watch?v=BkSdD5VtyRM" value={this.state.youtube_url} onChange={(evt) => this.applyYoutubeLink(evt.target.value)} />
                            </div>
                        </div>
                        <small>WARNING: Only include videos which you own the legal copyright to if you plan to monetize from the content. 
                            Including a video you do not legally own and monetizing it will result in your account being suspended.</small>
                    </p>
                </div>
                <div className="col-md-6" >
                    <div className="form-group">
                        <label>Video Title</label>
                        <input type="text" className="form-control" placeholder="Content title" value={this.state.campaign.title} onChange={(evt) => { let c = this.state.campaign; c.title = evt.target.value; this.setState({ campaign: c }) }} />
                    </div>
                    <div className="form-group">
                        <label>Who can view this ?</label>
                        <select className="form-control" value={this.state.campaign.subscription} onChange={(evt) => { let c = this.state.campaign; c.subscription = evt.target.value; this.setState({ campaign: c }) }}>
                            <option value="public">Everyone (Public )</option>
                        </select>
                        <p><small>Embedded content can only be public.</small></p>
                    </div>
                    <p>You must upload your file first !</p>
                    <button className="btn btn-block btn-primary" disabled={!this.state.valid_yt_link} onClick={this.createNewCampaign} ><i className="fa fa-check"></i> Publish video</button>
                </div>
            </div>
        )
    }
}

export default CreateYoutubeEmbedCampaign;