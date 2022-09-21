import React, { Component } from "react"
import { Link, Switch, Route } from 'react-router-dom'
import CreatorContentCard from "./CreatorContentCard"

class CreatorContentList extends Component{

    constructor(props) {
        super(props)
        this.state = {
            buyItems: 1,
            maxShowMobile: this.props.maxShowMobile ? this.props.maxShowMobile : 100,
        }
    }

    render() {
        return (
    <>
                <div class="row justify-content-center" >
                {this.props.content.map((item, idx) => {
                    return <CreatorContentCard redirect={this.props.redirect} classes={this.props.loadMore && idx > this.state.maxShowMobile ? 'd-none d-md-block' : '' } user={this.props.user} creator={this.props.creator} content={item} />
                })}
                </div>
                {this.props.loadMore ? 
                    <div class="row justify-content-center" >
                        <div class="col-lg-7 col-xs-12 col-md-7 order-md-1">
                            <div class="text-center">
                                <Link to={`/@${this.props.creator.username}/content`} class="btn btn-primary">Load More</Link>
                            </div>
                        </div>
                    </div>  : <></>}
    </>
        )
    }
}


export default CreatorContentList;