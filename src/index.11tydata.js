module.exports = {
  eleventyComputed: {
    noindex: function (data) {
      return data.pagination.pageNumber != 0;
    }
  }
};
