# Add/remove labels GitHub action

A tiny GitHub action that can add labels to pull requests and remove them from issues (and vice versa). Packaged as a tiny Golang binary in a Docker image for maximum performance.

## Usage

Use `urcomputeringpal/add-remove-labels-action` as a step in your workflow that processes `issues` and/or `pull_request` events.

Add labels:

```yaml
    - name: add-labels
      uses: urcomputeringpal/add-remove-labels-action@master
      with:
        GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"
        action: add
        labels: feature
```

Or remove them:

```yaml
    - name: add-labels
      uses: urcomputeringpal/add-remove-labels-action@master
      with:
        GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"
        action: remove
        labels: bug
```


## Advanced Usage 

This action pairs nicely with GitHub Actions' support for conditionally executing steps. In the following example, a label is added to or removed from a pull request based on the exit status of a theoretical `diff` utility.

```yaml
    strategy:
      fail-fast: false
      matrix:
        component:
          - app
          - assets
          - admin
    steps:
      # ...
      - name: diff
        id: diff
        run: |
          if ./pseudocode-diff-tool; then
              echo "::set-output name=changes::true"
          else
              echo "::set-output name=changes::false"
          fi
      - name: add-labels
        uses: urcomputeringpal/add-remove-labels-action@master
        if: "steps.diff.outputs.changes == 'true'"
        with:
          GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"
          action: add
          labels: "changes ${{ matrix.component }}"

      - name: remove-labels
        uses: urcomputeringpal/add-remove-labels-action@master
        if: "steps.diff.outputs.changes != 'true'"
        with:
          GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"
          action: remove
          labels: "changes ${{ matrix.component }}"
```
