import React, { Component } from 'react';
import { Grid, Navbar, Button, Form, FormGroup, InputGroup, FormControl, ListGroup, ListGroupItem, Glyphicon} from 'react-bootstrap';

class App extends Component {
  constructor(props) {
    super(props)
    this.state = {
      value: "",
      results: new Map()
    }
    this.handleChange = this.handleChange.bind(this)
    this.handleSearch = this.handleSearch.bind(this)
  }

  handleChange(e) {
    this.setState({ value: e.target.value })
  }

  handleSearch(e) {
    e.preventDefault()
    let query = this.state.value.trim()
    if (query.length === 0) {
      return
    }

    fetch(`https://pt-search-backend-4ywt76ki5q-uc.a.run.app/api/search?q=${query}`)
      .then((response) => response.json())
      .then((responseJson) => {
        this.setState({
          results: this._groupSearchResults(responseJson)
        })
      })
      .catch((error) => {

      })
  }

  _groupSearchResults(response) {
    let movieMap = new Map()
    for (let r of response) {
      let key = `${r.Title}:${r.Year}`
      if (!movieMap.has(key)) {
        movieMap.set(key, {
          title: r.Title,
          year: r.Year,
          resources: []
        })
      }
      movieMap.get(key).resources.push(r)
    }

    console.dir(movieMap)

    return movieMap
  }

  render() {
    return (
      <div>
        <Navbar>
          <Navbar.Header>
            <Navbar.Brand>
              <a href="/">PT-Search</a>
            </Navbar.Brand>
          </Navbar.Header>
        </Navbar>
        <Grid>
          <Form onSubmit={this.handleSearch}>
            <FormGroup>
              <InputGroup>
                <FormControl
                  type="text"
                  value={this.state.value}
                  placeholder="Type to search"
                  onChange={this.handleChange}
                />
                <InputGroup.Button>
                  <Button type="submit" onClick={this.handleSearch}>Search</Button>
                </InputGroup.Button>
              </InputGroup>
            </FormGroup>
          </Form>
          <ResultList results={this.state.results}></ResultList>
        </Grid>
      </div>
    );
  }
}

class MovieItem extends Component {
  constructor(props) {
    super(props);
    this.state = {
      downloadStatus: ""
    }
    this.handleDownload = this.handleDownload.bind(this)
  }

  handleDownload(e) {
    let apiUrl = `https://pt-search-backend-4ywt76ki5q-uc.a.run.app/api/queue?from=${this.props.from}&id=${this.props.id}`
    console.log('send download request: ' + apiUrl)

    this.setState({downloadStatus: "downloading"})
    fetch(apiUrl, {method: 'POST'})
      .then((response) => {
        if (!response.ok) {
          this.setState({downloadStatus: "failed"})
          throw new Error('failed to send download request')
        } else {
          this.setState({downloadStatus: "finished"})
        }
      }).catch((err) => {
        console.log(err)
      })
  }

  render() {
    const resourceItems = this.props.resources.map((r, i) => {
      return <ListGroupItem key={i}>
        <div>
          <a href={r.Tracker.url} target="_blank">
            {r.Tracker.title}
          </a>
          <p>{r.Tracker.sub_title}</p>
        </div>
        <div>
          <span>[{r.Tracker.from}] </span>
          <span>{r.Resource.Source} </span>
          <span>{r.Resource.Resolution} </span>
          <span>{r.Resource.Group} </span>
          <span>{r.Tracker.age} </span>
          <span>{r.Tracker.size} </span>
          <span>{r.Tracker.seeder}</span>
          <DownloadButton downloadStatus={this.state.downloadStatus} onClick={this.handleDownload}></DownloadButton>
        </div>
      </ListGroupItem>
    })

    return <div>
      <span>{this.props.title} {this.props.year}</span>
      <ListGroup>{resourceItems}</ListGroup>
    </div>
  }
}

const DownloadButton = (props) => {
  switch(props.downloadStatus) {
    case 'downloading':
      return (
        <Button bsStyle="link" bsSize="sm" onClick={props.onClick} disabled>
          <Glyphicon glyph="sort-by-attributes" />
        </Button>
      )
    case 'finished':
      return (
        <Button bsStyle="link" bsSize="sm" disabled>
          <Glyphicon glyph="saved" />
        </Button>
      )
    default :
      return (
        <Button bsStyle="link" bsSize="sm" onClick={props.onClick}>
          <Glyphicon glyph="download-alt"/>
        </Button>
      )
  }
}

const ResultList = (props) => {
  const movies = props.results
  const movieItems = []
  
  for (let m of movies.values()) {
    movieItems.push(
      <li key={m.title + '-' + m.year}>
        <MovieItem title={m.title} year={m.year} resources={m.resources}></MovieItem>
      </li>
    )
  }

  return <div>
    <ul>{movieItems}</ul>
  </div>
}

export default App;
