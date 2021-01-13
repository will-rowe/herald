This quick start document will walk through what to do when running **Herald** for the first time.

### Experiments

Experiments are equivalent to a MinKNOW experiment; i.e. where you record the information prior to sequencing.

**Herald** uses MinKNOW experiments to group samples.

#### Create an experiment

To create an experiment, first click on the `create experiment` button to access the submission form and then fill in the following fields:

---

`Experiment name`

This is the unique name to give this experiment. All data for this experiment (including any samples you create later) will be kept in a directory with this name.

> **note**: your value for `experiment name` can contain spaces but these will be replaced with underscores in all file names that are derived from this value.

---

`Output location`

This is where you want the experiment data to be stored.

---

Once you have entered both `Experiment name` and `Output location` and have clicked onto the next box, the form will autocomplete to indicate where the data will be stored.

> **example**: if you entered `my experiment` and `/Users/myname/data`, the form will autocomplete your `FAST5 data` directory to `/Users/myname/data/my_experiment/fast5_pass`

The app will also check the proposed directories - if `fast5_pass` or `fastq_pass` is found, the app will assume you are entering details for an existing experiment and will tag the sample as sequenced and basecalled already.

Fill in any other details and then click on the `create sample` button on the bottom of the form

A success message should appear and the number of available experiments will have increased in the app dashboard.

### Samples

### Processes

Once an experiment or sample has been tagged with a process, it can't be removed. The process will be marked as complete when it has finished.

note: If an existing experiment has been added (which already has fast5 or fastq data), the `sequence` and `basecall` tags will be added but marked to `complete` straight away.
