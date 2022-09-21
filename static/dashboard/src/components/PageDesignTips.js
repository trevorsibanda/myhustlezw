import React, {Component} from "react"
import ReactPlayer from "react-player";



class PageDesignTips extends Component {
    render() {
        return (
            <div className="box">
                <div className="box-body">
                    <h4 className="box-title">Tips for your page</h4>
                    <div className="player-wrapper">
                        <ReactPlayer width='100%' controls={true} height='100%' playing={false} url="https://www.youtube.com/watch?v=51Y83OTjmzg" className="react-player" />
                    </div>
                </div>
            </div>
        )
    }
}

export default PageDesignTips;