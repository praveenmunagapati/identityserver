- Make it possible that an organization adds a URL to a document to IYOSEE

- Make it accessible for the user
  - GET `https://itsyou.online/users/<userid>/see`: list all organizations that added something to this user's `see`
  - GET `https://itsyou.online/users/<userid>/see/<organization>` to list all documents added by this organization
  - POST `https://itsyou.online/users/<userid>/see/<organization>/<uniqueid>` to save a link
  - PUT `https://itsyou.online/users/<userid>/see/<organization>/<uniqueid>` to update a link

User can't add his own links

also need scope to allow organizations to write to this

- Store the following info of a document
  - link
  - category
  - uniqueid (name)
  - version
  - content_type
  - markdown short descr
  - markdown full descr
  - validity_date: start_date (/ end_date)
  - post_date: data that info was uploaded
  - required security level
  - signature (with priv key of user on rogerthat app) of all content above (e.g. json)
